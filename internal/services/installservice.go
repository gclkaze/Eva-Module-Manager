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
func (inst *InstallService) InstallAllFromPath(ctx context.Context, token string, p string) error {
	p = strings.TrimSpace(p)
	if p == "" {
		return fmt.Errorf("project file path provided but it is empty")
	}
	eva := inst.bookKeepingService.DefaultEvaFileName()
	if !utils.FolderExists(p) {
		return fmt.Errorf("path '%s' is not a valid folder", p)
	}

	filename := filepath.Join(p, eva)
	if !utils.FileExists(filename) {
		return fmt.Errorf("project file '%s' does not exist", filename)
	}

	return inst.InstallAllFromProjectFile(ctx, token, filename)
}

func (inst *InstallService) InstallModuleVersion(ctx context.Context, token string, module string, version string) error {
	return nil
}

// with ./eva.json
func (inst *InstallService) InstallAllFromProject(ctx context.Context, token string) error {
	eva := inst.bookKeepingService.DefaultEvaFileName()
	filename := filepath.Join(inst.cwd, eva)
	if !utils.FileExists(filename) {
		return fmt.Errorf("project file '%s' does not exist", filename)
	}

	return inst.InstallAllFromProjectFile(ctx, token, filename)
}

func (inst *InstallService) InstallAllFromProjectFile(ctx context.Context, token string, evafile string) error {
	absPath, err := inst.bookKeepingService.VerifyExisting(evafile)
	if err != nil {
		return err
	}

	inst.output.VerboseInfo(fmt.Sprintf("Project file at '%s' was verified successfully.", evafile))
	err = inst.bookKeepingService.LoadExisting(absPath)
	if err != nil {
		return err
	}

	theProject := inst.bookKeepingService.Get()
	//we got the modules, lets install the ones not there->the ones not in our file system

	modulesFolder := inst.bookKeepingService.DefaultEvaModulesFolder()
	modulesFolder = filepath.Join(inst.cwd, modulesFolder)

	if !utils.FolderExists(modulesFolder) {
		inst.output.VerboseInfo(fmt.Sprintf("Creating moduesl folder '%s'.", modulesFolder))
		err = utils.CreateFolder(modulesFolder)
		if err != nil {
			return err
		}
		inst.output.VerboseInfo("Attempting a clean installation...")
		return inst.CleanInstallAllFromProjectFile(ctx, token, theProject, modulesFolder)
	}
	inst.output.VerboseInfo("Attempting installation of the non-existing modules..")
	err = inst.DirtyInstallAllFromProjectFile(ctx, token, theProject, modulesFolder)
	return err
}

func (inst *InstallService) DirtyInstallAllFromProjectFile(ctx context.Context, token string, theProject eva.EvaProject, saveLocation string) error {
	for key, info := range theProject.Modules {
		moduleName := key
		theModule := info

		if !utils.FolderIsEmpty(theModule.InstallationFolder) {
			inst.output.VerboseInfo(fmt.Sprintf("Skipping download & installation of module '%s'.", moduleName))
			continue
		}

		module, version, err := utils.ParseModuleReleaseVersion(moduleName)
		if err != nil {
			return err
		}

		modulePath := filepath.Join(saveLocation, module)
		if !utils.FolderExists(modulePath) {
			err = utils.CreateFolder(modulePath)
			if err != nil {
				inst.output.Error(err)
				return fmt.Errorf("coulnd't create module folder %s", modulePath)
			}
		}

		modulePath = filepath.Join(modulePath, version)
		if !utils.FolderExists(modulePath) {
			err = utils.CreateFolder(modulePath)
			if err != nil {
				inst.output.Error(err)
				return fmt.Errorf("coulnd't create module version folder %s", modulePath)
			}
		}
		inst.output.VerboseInfo(fmt.Sprintf("Downloading '%s@%s' and storing it at %s.", module, version, modulePath))
		err = inst.releaseService.DownloadRelease(ctx, token, module, version, modulePath)
		if err != nil {
			return err
		}
		inst.output.VerboseInfo(fmt.Sprintf("Download completed of '%s@%s', storing it at %s.", module, version, modulePath))
		err = inst.BuildModuleFolderFromTar(ctx, module, version, modulePath)
		if err != nil {
			inst.output.Error(err)
			return nil
		}
	}
	return nil
}

func (inst *InstallService) CleanInstallAllFromProjectFile(ctx context.Context, token string, theProject eva.EvaProject, saveLocation string) error {
	for key := range theProject.Modules {
		moduleName := key
		//theModule := info

		module, version, err := utils.ParseModuleReleaseVersion(moduleName)
		if err != nil {
			return err
		}

		modulePath := filepath.Join(saveLocation, module)
		if !utils.FolderExists(modulePath) {
			err = utils.CreateFolder(modulePath)
			if err != nil {
				inst.output.Error(err)
				return fmt.Errorf("coulnd't create module folder %s", modulePath)
			}
		}

		modulePath = filepath.Join(modulePath, version)
		if !utils.FolderExists(modulePath) {
			err = utils.CreateFolder(modulePath)
			if err != nil {
				inst.output.Error(err)
				return fmt.Errorf("coulnd't create module version folder %s", modulePath)
			}
		}
		inst.output.VerboseInfo(fmt.Sprintf("Downloading '%s@%s' and storing it at %s.", module, version, modulePath))
		err = inst.releaseService.DownloadRelease(ctx, token, module, version, modulePath)
		if err != nil {
			return err
		}
		inst.output.VerboseInfo(fmt.Sprintf("Download completed of '%s@%s', storing it at %s.", module, version, modulePath))
		err = inst.BuildModuleFolderFromTar(ctx, module, version, modulePath)
		if err != nil {
			inst.output.Error(err)
			return nil
		}
	}
	return nil
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
