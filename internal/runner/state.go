package runner

import (
	"encoding/json"
	"github.com/datadog/stratus-red-team/internal/utils"
	"github.com/datadog/stratus-red-team/pkg/stratus"
	"log"
	"os"
	"path"
	"path/filepath"
)

const StratusStateDirectoryName = ".stratus-red-team"

type StateManager struct {
	RootDirectory string
	Technique     *stratus.AttackTechnique
	FileSystem    FileSystem
}

type FileSystem interface {
	FileExists(string) bool
	CreateDirectory(string, os.FileMode) error
	WriteFile(string, []byte, os.FileMode) error
	ReadFile(string) ([]byte, error)
}

type LocalFileSystem struct{}

func (m *LocalFileSystem) FileExists(fileName string) bool {
	return utils.FileExists(fileName)
}

func (m *LocalFileSystem) CreateDirectory(dir string, mode os.FileMode) error {
	return os.Mkdir(dir, mode)
}

func (m *LocalFileSystem) WriteFile(file string, content []byte, mode os.FileMode) error {
	return os.WriteFile(file, content, mode)
}

func (m *LocalFileSystem) ReadFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func NewStateManager(technique *stratus.AttackTechnique) *StateManager {
	homeDirectory, _ := os.UserHomeDir()
	stateManager := StateManager{
		RootDirectory: filepath.Join(homeDirectory, StratusStateDirectoryName),
		Technique:     technique,
		FileSystem:    &LocalFileSystem{},
	}
	stateManager.initialize()
	return &stateManager
}

func (m *StateManager) initialize() {
	if !m.FileSystem.FileExists(m.RootDirectory) {
		log.Println("Creating " + m.RootDirectory + " as it doesn't exist yet")
		err := m.FileSystem.CreateDirectory(m.RootDirectory, 0744)
		if err != nil {
			panic("Unable to create persistent directory: " + err.Error())
		}
	}
}

func (m *StateManager) ExtractTechniqueTerraformFile() error {
	terraformDirectory := filepath.Join(m.RootDirectory, m.Technique.ID)
	terraformFile := filepath.Join(terraformDirectory, "main.tf")

	if m.FileSystem.FileExists(terraformDirectory) {
		return nil
	}
	err := m.FileSystem.CreateDirectory(terraformDirectory, 0744)
	if err != nil {
		return err
	}
	return m.FileSystem.WriteFile(terraformFile, m.Technique.PrerequisitesTerraformCode, 0644)
}

func (m *StateManager) GetTechniqueOutputs() (map[string]string, error) {
	outputPath := path.Join(m.RootDirectory, m.Technique.ID, ".terraform-outputs")
	outputs := make(map[string]string)

	// If we have persisted Terraform outputs on disk, read them
	if m.FileSystem.FileExists(outputPath) {
		outputString, err := m.FileSystem.ReadFile(outputPath)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(outputString, &outputs)
		if err != nil {
			return nil, err
		}
	}
	return outputs, nil
}

func (m *StateManager) WriteTerraformOutputs(outputs map[string]string) error {
	outputPath := path.Join(m.RootDirectory, m.Technique.ID, ".terraform-outputs")
	outputString, err := json.Marshal(outputs)
	if err != nil {
		return err
	}
	return m.FileSystem.WriteFile(outputPath, outputString, 0744)
}

func (m *StateManager) GetTechniqueState() AttackTechniqueState {
	rawState, _ := m.FileSystem.ReadFile(filepath.Join(m.RootDirectory, m.Technique.ID, ".state"))
	return AttackTechniqueState(rawState)
}

func (m *StateManager) SetTechniqueState(state AttackTechniqueState) error {
	file := filepath.Join(m.RootDirectory, m.Technique.ID, ".state")
	return m.FileSystem.WriteFile(file, []byte(state), 0744)
}

func (m *StateManager) GetRootDirectory() string {
	return m.RootDirectory
}
