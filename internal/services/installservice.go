package services

import (
	"context"
	"emm/internal/backend"
	"emm/internal/models"
	"emm/internal/models/eva"
	"emm/internal/output"
	"emm/pkg/utils"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"
)

type InstallService struct {
	bookKeepingService *ProjectBookkeepingService
	releaseService     *ModuleReleaseService
	props              *properties.Properties
	output             output.Printer
	currentUser        *models.User
	backend            *backend.Backend

	cwd string
}

func NewInstallService(cwd string, bookkeeping *ProjectBookkeepingService, releaseService *ModuleReleaseService) *InstallService {
	return &InstallService{cwd: cwd, bookKeepingService: bookkeeping, releaseService: releaseService}
}

func (inst *InstallService) SetBackend(backend *backend.Backend) {
	inst.backend = backend
}

func (inst *InstallService) SetPrinter(output output.Printer) {
	inst.output = output
}

func (inst *InstallService) SetProperties(props *properties.Properties) {
	inst.props = props
	//inst.currentUserKey = props.GetString("currentUserKey", "EMMCurrentUser")
}
func (inst *InstallService) InstallAllFromPath(ctx context.Context, token string, p string) (*models.InstallationSummary, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return nil, fmt.Errorf("project file path provided but it is empty")
	}
	eva := inst.bookKeepingService.DefaultEvaFileName()
	if !utils.FolderExists(p) {
		return nil, fmt.Errorf("path '%s' is not a valid folder", p)
	}

	filename := filepath.Join(p, eva)
	if !utils.FileExists(filename) {
		return nil, fmt.Errorf("project file '%s' does not exist", filename)
	}

	return inst.InstallAllFromProjectFile(ctx, token, filename)
}

func (inst *InstallService) FindRelease(module string, theProject *eva.EvaProject) (string, string, error) {
	for key := range theProject.Modules {
		theModule, theVersion, err := utils.ParseModuleReleaseVersion(key)
		if err != nil {
			return "", "", err
		}
		if theModule == module {
			return theModule, theVersion, nil
		}
	}
	return "", "", nil
}

func (inst *InstallService) InstallModuleVersion(ctx context.Context, token string, module string, version string, path *string) (*models.InstallationSummary, error) {
	theProject, absPath, err := inst.GetEvaProjectFile(path)
	if err != nil {
		return nil, err
	}

	theModule, theVersion, err := inst.FindRelease(module, theProject)
	if err != nil {
		return nil, err
	}

	modulesFolder := ""
	if theProject.ModulesFolder == "" {
		modulesFolder = inst.bookKeepingService.DefaultEvaModulesFolder()
		modulesFolder = filepath.Join(inst.cwd, modulesFolder)
	} else {
		if filepath.IsAbs(theProject.ModulesFolder) {
			modulesFolder = theProject.ModulesFolder
		} else {
			dir := filepath.Dir(absPath)
			modulesFolder = filepath.Join(dir, theProject.ModulesFolder)
		}
	}
	summary := models.NewInstallationSummary()
	//we need to see if we know this module.
	if theModule == "" && theVersion == "" {
		//the module wasn't found
		/*installed*/
		//if we do not, we download it, install it and add it in eva.json
		_, err = inst.InstallNewModule(ctx, token, module, version, summary, modulesFolder, fmt.Sprintf("%s@%s", module, version))
		return summary, err
	}
	//if we do, we check the versions are equal, then we skip
	//if there are not, we inform the user that we will update the module and its version, then we download, we install and then add it in eva.json
	//if there is no version, we fetch the latest. then we check again to see if we know it. if the input and stored versions are identical, we print a message and exit
	_, err = inst.InstallKnownModule(ctx, token, theModule, theVersion, summary, modulesFolder, module, version)
	return summary, err
}

func (inst *InstallService) HandleMissingVersions(ctx context.Context, token string, module string, storedVersion string, inputVersion string, summary *models.InstallationSummary, saveLocation string, originalModuleName string) (bool, error) {
	storedMissesVersion := inst.isLatestVersion(storedVersion)
	inputMissesVersion := inst.isLatestVersion(inputVersion)

	inputModuleName := fmt.Sprintf("%s@%s", module, inputVersion)
	//if the stored is missing the version, this means that the user may have changed it by hand.
	// what we will do is to fetch the latest version and install it. and thats it
	if storedMissesVersion && inputMissesVersion {
		inst.output.VerboseWarn(fmt.Sprintf("stored module name is '%s' and input module name is '%s', both miss version", originalModuleName, inputModuleName))
		version, err := inst.handleLatestVersion(ctx, token, module)
		if err != nil {
			return false, err
		}
		//installation as usual here & update the project file
		return inst.InstallNewModuleWithKnownVersions(ctx, token, module, version, summary, saveLocation, originalModuleName)

	} else if storedMissesVersion && !inputMissesVersion {
		inst.output.VerboseWarn(fmt.Sprintf("stored module name is '%s' and input module name is '%s', stored only misses version", originalModuleName, inputModuleName))
		//we take into account the inputVersion
		return inst.InstallNewModuleWithKnownVersions(ctx, token, module, inputVersion, summary, saveLocation, originalModuleName)

	} else if !storedMissesVersion && inputMissesVersion {
		inst.output.VerboseWarn(fmt.Sprintf("stored module name is '%s' and input module name is '%s', new only misses version", originalModuleName, inputModuleName))
		//the inputVersion misses the exact version, we fetch the latest
		version, err := inst.handleLatestVersion(ctx, token, module)
		if err != nil {
			return false, err
		}

		//installation as usual here & update the project file
		return inst.InstallNewModuleWithKnownVersions(ctx, token, module, version, summary, saveLocation, originalModuleName)
	}
	inst.output.VerboseWarn(fmt.Sprintf("stored module name is '%s' and input module name is '%s', none misses version", originalModuleName, inputModuleName))
	//all versions are present
	return false, nil

}

func (inst *InstallService) InstallKnownModule(ctx context.Context, token string, module string, version string, summary *models.InstallationSummary, saveLocation string, inputModule string, inputVersion string) (bool, error) {
	originalModuleName := fmt.Sprintf("%s@%s", module, version)
	inputModuleName := fmt.Sprintf("%s@%s", inputModule, inputVersion)

	installed, err := inst.HandleMissingVersions(ctx, token, module, version, inputVersion, summary, saveLocation, originalModuleName)
	if err != nil {
		inst.output.VerboseWarn(fmt.Sprintf("failed to handle missing versions where stored module name is '%s' and input module name is '%s'", originalModuleName, inputModuleName))
		return false, err
	}

	if installed {
		return true, nil
	}

	//if we do, we check the versions are equal, then we skip
	if utils.ReleasesAreEquivalent(module, version, inputModule, inputVersion) {
		inst.output.VerboseInfo(fmt.Sprintf("Skipping download & installation of module '%s'.", originalModuleName))
		summary.Skipped += 1
		return false, nil
	}

	//if there are not, we inform the user that we will update the module and its version, then we download, we install and then add it in eva.json
	//lets inform the user about the inputModule@inputVersion
	inst.output.Info(fmt.Sprintf("the currently installed version of the module '%s' will be replaced with '%s'.", version, inputModuleName))

	modulePath, err := inst.incrementalFolderCreation(saveLocation, inputModule, inputVersion, inputModuleName, summary)
	if err != nil {
		return false, err
	}
	if modulePath == "" {
		//it was skipped
		err = inst.bookKeepingService.ReplaceAndCommitModule(originalModuleName, inputModule, inputVersion)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	inst.output.VerboseInfo(fmt.Sprintf("Downloading '%s@%s' and storing it at %s.", inputModule, inputVersion, modulePath))
	err = inst.releaseService.DownloadRelease(ctx, token, inputModule, inputVersion, modulePath)
	if err != nil {
		summary.Failed += 1
		utils.CleanAndRemoveFolder(modulePath)
		return false, err
	}
	inst.output.VerboseInfo(fmt.Sprintf("Download completed of '%s@%s', storing it at %s.", inputModule, inputVersion, modulePath))
	err = inst.BuildModuleFolderFromTar(ctx, inputModule, inputVersion, modulePath)
	if err != nil {
		summary.Failed += 1
		inst.output.Error(err)
		return false, nil
	}
	summary.ProcessedCounter += 1
	summary.Success += 1

	err = inst.bookKeepingService.ReplaceAndCommitModule(originalModuleName, inputModule, inputVersion)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (inst *InstallService) incrementalFolderCreation(saveLocation string, module string, version string, originalModuleName string, summary *models.InstallationSummary) (string, error) {
	var err error
	modulePath := filepath.Join(saveLocation, module)
	if !utils.FolderExists(modulePath) {
		err = utils.CreateFolder(modulePath)
		if err != nil {
			summary.Failed += 1
			inst.output.Error(err)
			return "", fmt.Errorf("coulnd't create module folder %s", modulePath)
		}
	}

	modulePath = filepath.Join(modulePath, version)
	if !utils.FolderExists(modulePath) {
		err = utils.CreateFolder(modulePath)
		if err != nil {
			summary.Failed += 1
			inst.output.Error(err)
			return "", fmt.Errorf("coulnd't create module version folder %s", modulePath)
		}
	}

	if !utils.FolderIsEmpty(modulePath) {
		inst.output.VerboseInfo(fmt.Sprintf("Skipping download & installation of module '%s'.", originalModuleName))
		summary.Skipped += 1
		return "", nil
	}
	return modulePath, nil
}

func (inst InstallService) isLatestVersion(version string) bool {
	return version == "" || version == "latest"
}

func (inst *InstallService) InstallNewModule(ctx context.Context, token string, module string, version string, summary *models.InstallationSummary, saveLocation string, originalModuleName string) (bool, error) {
	var err error
	if inst.isLatestVersion(version) {
		version, err = inst.handleLatestVersion(ctx, token, module)
		if err != nil {
			return false, err
		}
	}

	modulePath, err := inst.incrementalFolderCreation(saveLocation, module, version, originalModuleName, summary)
	if err != nil {
		return false, err
	}
	if modulePath == "" {
		//it was skipped
		err = inst.bookKeepingService.AddAndCommitModule(module, version)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	inst.output.VerboseInfo(fmt.Sprintf("Downloading '%s@%s' and storing it at %s.", module, version, modulePath))
	err = inst.releaseService.DownloadRelease(ctx, token, module, version, modulePath)
	if err != nil {
		summary.Failed += 1
		utils.CleanAndRemoveFolder(modulePath)
		return false, err
	}
	inst.output.VerboseInfo(fmt.Sprintf("Download completed of '%s@%s', storing it at %s.", module, version, modulePath))
	err = inst.BuildModuleFolderFromTar(ctx, module, version, modulePath)
	if err != nil {
		summary.Failed += 1
		inst.output.Error(err)
		return false, nil
	}
	summary.ProcessedCounter += 1
	summary.Success += 1

	err = inst.bookKeepingService.AddAndCommitModule(module, version)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (inst *InstallService) InstallNewModuleWithKnownVersions(ctx context.Context, token string, module string, version string, summary *models.InstallationSummary, saveLocation string, originalModuleName string) (bool, error) {
	var err error
	modulePath, err := inst.incrementalFolderCreation(saveLocation, module, version, originalModuleName, summary)
	if err != nil {
		return false, err
	}
	if modulePath == "" {
		//it was skipped
		err = inst.bookKeepingService.EnsureModuleAndCommit(module, version)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	inst.output.VerboseInfo(fmt.Sprintf("Downloading '%s@%s' and storing it at %s.", module, version, modulePath))
	err = inst.releaseService.DownloadRelease(ctx, token, module, version, modulePath)
	if err != nil {
		summary.Failed += 1
		utils.CleanAndRemoveFolder(modulePath)
		return false, err
	}
	inst.output.VerboseInfo(fmt.Sprintf("Download completed of '%s@%s', storing it at %s.", module, version, modulePath))
	err = inst.BuildModuleFolderFromTar(ctx, module, version, modulePath)
	if err != nil {
		summary.Failed += 1
		inst.output.Error(err)
		return false, nil
	}
	summary.ProcessedCounter += 1
	summary.Success += 1

	err = inst.bookKeepingService.EnsureModuleAndCommit(module, version)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (inst *InstallService) GetEvaProjectFile(path *string) (*eva.EvaProject, string, error) {
	eva := inst.bookKeepingService.DefaultEvaFileName()
	evafile := ""
	if path == nil {
		evafile = filepath.Join(inst.cwd, eva)
	} else {
		evafile = filepath.Join(*path, eva)
	}

	if !utils.FileExists(evafile) {
		return nil, "", fmt.Errorf("project file '%s' does not exist", evafile)
	}
	absPath, err := inst.bookKeepingService.VerifyExisting(evafile)
	if err != nil {
		return nil, "", err
	}

	inst.output.VerboseInfo(fmt.Sprintf("Project file at '%s' was verified successfully.", evafile))
	err = inst.bookKeepingService.LoadExisting(absPath)
	if err != nil {
		return nil, "", err
	}

	theProject := inst.bookKeepingService.Get()
	return &theProject, absPath, nil
}

// with ./eva.json
func (inst *InstallService) InstallAllFromProject(ctx context.Context, token string) (*models.InstallationSummary, error) {
	eva := inst.bookKeepingService.DefaultEvaFileName()
	filename := filepath.Join(inst.cwd, eva)
	if !utils.FileExists(filename) {
		return nil, fmt.Errorf("project file '%s' does not exist", filename)
	}

	return inst.InstallAllFromProjectFile(ctx, token, filename)
}

func (inst *InstallService) InstallAllFromProjectFile(ctx context.Context, token string, evafile string) (*models.InstallationSummary, error) {
	absPath, err := inst.bookKeepingService.VerifyExisting(evafile)
	if err != nil {
		return nil, err
	}

	inst.output.VerboseInfo(fmt.Sprintf("Project file at '%s' was verified successfully.", evafile))
	err = inst.bookKeepingService.LoadExisting(absPath)
	if err != nil {
		return nil, err
	}

	theProject := inst.bookKeepingService.Get()
	//we got the modules, lets install the ones not there->the ones not in our file system

	//We respect the folder described in eva.json
	modulesFolder := ""
	if theProject.ModulesFolder == "" {
		modulesFolder = inst.bookKeepingService.DefaultEvaModulesFolder()
		modulesFolder = filepath.Join(inst.cwd, modulesFolder)
	} else {
		if filepath.IsAbs(theProject.ModulesFolder) {
			modulesFolder = theProject.ModulesFolder
		} else {
			dir := filepath.Dir(absPath)
			modulesFolder = filepath.Join(dir, theProject.ModulesFolder)
		}
	}

	if !utils.FolderExists(modulesFolder) {
		inst.output.VerboseInfo(fmt.Sprintf("Creating moduesl folder '%s'.", modulesFolder))
		err = utils.CreateFolder(modulesFolder)
		if err != nil {
			return nil, err
		}
		inst.output.VerboseInfo("Attempting a clean installation...")
		return inst.CleanInstallAllFromProjectFile(ctx, token, theProject, modulesFolder)
	}
	inst.output.VerboseInfo("Attempting installation of the non-existing modules..")
	return inst.DirtyInstallAllFromProjectFile(ctx, token, theProject, modulesFolder)
}

func (inst *InstallService) DirtyInstallAllFromProjectFile(ctx context.Context, token string, theProject eva.EvaProject, saveLocation string) (*models.InstallationSummary, error) {
	summary := models.NewInstallationSummary()
	summary.Total = len(theProject.Modules)

	for key := range theProject.Modules {
		moduleName := key
		//theModule := info

		module, version, err := utils.ParseModuleReleaseVersion(moduleName)
		if err != nil {
			summary.Failed += 1
			return summary, err
		}
		if inst.isLatestVersion(version) {
			version, err = inst.handleLatestVersion(ctx, token, module)
			if err != nil {
				return summary, err
			}
		}

		modulePath := filepath.Join(saveLocation, module)
		if !utils.FolderExists(modulePath) {
			err = utils.CreateFolder(modulePath)
			if err != nil {
				summary.Failed += 1
				inst.output.Error(err)
				return summary, fmt.Errorf("coulnd't create module folder %s", modulePath)
			}
		}

		modulePath = filepath.Join(modulePath, version)
		if !utils.FolderExists(modulePath) {
			err = utils.CreateFolder(modulePath)
			if err != nil {
				summary.Failed += 1
				inst.output.Error(err)
				return summary, fmt.Errorf("coulnd't create module version folder %s", modulePath)
			}
		}

		if !utils.FolderIsEmpty(modulePath) {
			inst.output.VerboseInfo(fmt.Sprintf("Skipping download & installation of module '%s'.", moduleName))
			summary.Skipped += 1
			continue
		}

		inst.output.VerboseInfo(fmt.Sprintf("Downloading '%s@%s' and storing it at %s.", module, version, modulePath))
		err = inst.releaseService.DownloadRelease(ctx, token, module, version, modulePath)
		if err != nil {
			summary.Failed += 1
			return summary, err
		}
		inst.output.VerboseInfo(fmt.Sprintf("Download completed of '%s@%s', storing it at %s.", module, version, modulePath))
		err = inst.BuildModuleFolderFromTar(ctx, module, version, modulePath)
		if err != nil {
			summary.Failed += 1
			inst.output.Error(err)
			return summary, nil
		}
		summary.ProcessedCounter += 1
		summary.Success += 1
	}

	inst.bookKeepingService.Save()
	return summary, nil
}

func (inst *InstallService) handleLatestVersion(ctx context.Context, token string, module string) (string, error) {
	//version = "latest"
	//we need the latest release number
	theRelease, err := inst.releaseService.GetLatestModuleRelease(ctx, token, module)
	if err != nil {
		inst.output.Error(err)
		return "", fmt.Errorf("module %s has no latest release", module)
	}

	if theRelease == nil {
		inst.output.Error(err)
		return "", fmt.Errorf("module %s has no latest release", module)
	}

	version := theRelease.Version
	inst.output.VerboseInfo(fmt.Sprintf("Resolved latest version for '%s@%s'. Updating project file .", module, version))
	err = inst.bookKeepingService.ResolveAndPin(module, version)
	if err != nil {
		return "", fmt.Errorf("couldn't resolve latest version for : '%s@%s'", module, version)
	}
	return version, nil
}

func (inst *InstallService) CleanInstallAllFromProjectFile(ctx context.Context, token string, theProject eva.EvaProject, saveLocation string) (*models.InstallationSummary, error) {
	summary := models.NewInstallationSummary()
	summary.Total = len(theProject.Modules)

	for key := range theProject.Modules {
		moduleName := key
		//theModule := info

		module, version, err := utils.ParseModuleReleaseVersion(moduleName)
		if err != nil {
			summary.Failed += 1
			return summary, err
		}
		if inst.isLatestVersion(version) {
			version, err = inst.handleLatestVersion(ctx, token, module)
			if err != nil {
				return summary, err
			}
		}

		modulePath := filepath.Join(saveLocation, module)
		if !utils.FolderExists(modulePath) {
			err = utils.CreateFolder(modulePath)
			if err != nil {
				inst.output.Error(err)
				summary.Failed += 1
				return summary, fmt.Errorf("coulnd't create module folder %s", modulePath)
			}
		}

		modulePath = filepath.Join(modulePath, version)

		if !utils.FolderExists(modulePath) {
			err = utils.CreateFolder(modulePath)
			if err != nil {
				summary.Failed += 1
				inst.output.Error(err)
				return summary, fmt.Errorf("coulnd't create module version folder %s", modulePath)
			}
		}
		inst.output.VerboseInfo(fmt.Sprintf("Downloading '%s@%s' and storing it at %s.", module, version, modulePath))
		err = inst.releaseService.DownloadRelease(ctx, token, module, version, modulePath)
		if err != nil {
			summary.Failed += 1
			return summary, err
		}
		inst.output.VerboseInfo(fmt.Sprintf("Download completed of '%s@%s', storing it at %s.", module, version, modulePath))
		err = inst.BuildModuleFolderFromTar(ctx, module, version, modulePath)
		if err != nil {
			summary.Failed += 1
			inst.output.Error(err)
			return summary, nil
		}

		summary.ProcessedCounter += 1
		summary.Success += 1
	}

	return summary, nil
}

func (inst *InstallService) BuildModuleFolderFromTar(ctx context.Context, module string, version string, distLocation string) error {
	if !utils.FolderExists(distLocation) {
		return fmt.Errorf("base module folder does not exist: '%s'", distLocation)
	}

	theFile := inst.releaseService.GetModuleTarBallName(module, version)
	theFile = filepath.Join(distLocation, theFile)
	if !utils.FileExists(theFile) {
		return fmt.Errorf("module file '%s' is absent ", theFile)
	}
	inst.output.VerboseInfo(fmt.Sprintf("Unziping tar ball '%s' and storing at '%s'.", theFile, distLocation))
	err := utils.UntarGzToDir(theFile, distLocation)
	if err != nil {
		inst.output.Error(fmt.Errorf("couldn't unzip module file : '%s'", theFile))
		return err
	}
	return utils.DeleteFile(theFile)
}
