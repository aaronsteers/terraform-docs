package module

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadModuleWithOptions(t *testing.T) {
	assert := assert.New(t)

	options, _ := NewOptions().With(&Options{
		Path: filepath.Join("testdata", "full-example"),
	})
	module, err := LoadWithOptions(options)

	assert.Nil(err)
	assert.Equal(true, module.HasHeader())
	assert.Equal(true, module.HasInputs())
	assert.Equal(true, module.HasOutputs())
	assert.Equal(true, module.HasProviders())
}

func TestLoadModule(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "load module from path",
			path:    "full-example",
			wantErr: false,
		},
		{
			name:    "load module from path",
			path:    "non-exist",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			_, err := loadModule(filepath.Join("testdata", tt.path))
			if tt.wantErr {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
			}
		})
	}
}

func TestLoadHeader(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "load module header from path",
			path:     "full-example",
			expected: "Example of 'foo_bar' module in `foo_bar.tf`.\n\n- list item 1\n- list item 2\n\nEven inline **formatting** in _here_ is possible.\nand some [link](https://domain.com/)",
		},
		{
			name:     "load module header from path",
			path:     "empty-header",
			expected: "",
		},
		{
			name:     "load module header from path",
			path:     "non-exist",
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			actual := loadHeader(filepath.Join("testdata", tt.path))
			assert.Equal(tt.expected, actual)
		})
	}
}

func TestLoadInputs(t *testing.T) {
	type expected struct {
		inputs    int
		requireds int
		optionals int
	}
	tests := []struct {
		name     string
		path     string
		expected expected
	}{
		{
			name: "load module inputs from path",
			path: "full-example",
			expected: expected{
				inputs:    6,
				requireds: 2,
				optionals: 4,
			},
		},
		{
			name: "load module inputs from path",
			path: "no-required-inputs",
			expected: expected{
				inputs:    6,
				requireds: 0,
				optionals: 6,
			},
		},
		{
			name: "load module inputs from path",
			path: "no-optional-inputs",
			expected: expected{
				inputs:    6,
				requireds: 6,
				optionals: 0,
			},
		},
		{
			name: "load module inputs from path",
			path: "no-inputs",
			expected: expected{
				inputs:    0,
				requireds: 0,
				optionals: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			module, _ := loadModule(filepath.Join("testdata", tt.path))
			inputs, requireds, optionals := loadInputs(module)

			assert.Equal(tt.expected.inputs, len(inputs))
			assert.Equal(tt.expected.requireds, len(requireds))
			assert.Equal(tt.expected.optionals, len(optionals))
		})
	}
}

func TestLoadOutputs(t *testing.T) {
	type expected struct {
		outputs int
	}
	tests := []struct {
		name     string
		path     string
		expected expected
	}{
		{
			name: "load module outputs from path",
			path: "full-example",
			expected: expected{
				outputs: 3,
			},
		},
		{
			name: "load module outputs from path",
			path: "no-outputs",
			expected: expected{
				outputs: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			options := NewOptions()
			module, _ := loadModule(filepath.Join("testdata", tt.path))
			outputs, err := loadOutputs(module, options)

			assert.Nil(err)
			assert.Equal(tt.expected.outputs, len(outputs))

			for _, v := range outputs {
				assert.Equal(false, v.ShowValue)
			}
		})
	}
}

func TestLoadOutputsValues(t *testing.T) {
	type expected struct {
		outputs int
	}
	tests := []struct {
		name       string
		path       string
		outputPath string
		expected   expected
		wantErr    bool
	}{
		{
			name:       "load module outputs with values from path",
			path:       "full-example",
			outputPath: "output-values.json",
			expected: expected{
				outputs: 3,
			},
			wantErr: false,
		},
		{
			name:       "load module outputs with values from path",
			path:       "full-example",
			outputPath: "no-file.json",
			expected:   expected{},
			wantErr:    true,
		},
		{
			name:       "load module outputs with values from path",
			path:       "no-outputs",
			outputPath: "no-file.json",
			expected:   expected{},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			options, _ := NewOptions().With(&Options{
				OutputValues:     true,
				OutputValuesPath: filepath.Join("testdata", tt.path, tt.outputPath),
			})
			module, _ := loadModule(filepath.Join("testdata", tt.path))
			outputs, err := loadOutputs(module, options)

			if tt.wantErr {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.Equal(tt.expected.outputs, len(outputs))

				for _, v := range outputs {
					assert.Equal(true, v.ShowValue)
				}
			}
		})
	}
}

func TestLoadProviders(t *testing.T) {
	type expected struct {
		providers int
	}
	tests := []struct {
		name     string
		path     string
		expected expected
	}{
		{
			name: "load module providers from path",
			path: "full-example",
			expected: expected{
				providers: 3,
			},
		},
		{
			name: "load module providers from path",
			path: "no-providers",
			expected: expected{
				providers: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			module, _ := loadModule(filepath.Join("testdata", tt.path))
			providers := loadProviders(module)

			assert.Equal(tt.expected.providers, len(providers))
		})
	}
}

func TestLoadComments(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		fileName   string
		lineNumber int
		expected   string
	}{
		{
			name:       "load resource comment from file",
			path:       "full-example",
			fileName:   "variables.tf",
			lineNumber: 2,
			expected:   "D description",
		},
		{
			name:       "load resource comment from file",
			path:       "full-example",
			fileName:   "variables.tf",
			lineNumber: 16,
			expected:   "A Description in multiple lines",
		},
		{
			name:       "load resource comment from file with wrong line number",
			path:       "full-example",
			fileName:   "variables.tf",
			lineNumber: 100,
			expected:   "",
		},
		{
			name:       "load resource comment from non existing file",
			path:       "full-example",
			fileName:   "non-exist.tf",
			lineNumber: 5,
			expected:   "",
		},
		{
			name:       "load resource comment from non existing file",
			path:       "non-exist",
			fileName:   "variables.tf",
			lineNumber: 5,
			expected:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			actual := loadComments(filepath.Join("testdata", tt.path, tt.fileName), tt.lineNumber)
			assert.Equal(tt.expected, actual)
		})
	}
}

func TestSortItems(t *testing.T) {
	type expected struct {
		inputs    []string
		required  []string
		optional  []string
		outputs   []string
		providers []string
	}
	tests := []struct {
		name     string
		path     string
		sort     *SortBy
		expected expected
	}{
		{
			name: "sort module items",
			path: "full-example",
			sort: &SortBy{Name: false, Required: false},
			expected: expected{
				inputs:    []string{"D", "B", "E", "A", "C", "F"},
				required:  []string{"A", "F"},
				optional:  []string{"D", "B", "E", "C"},
				outputs:   []string{"C", "A", "B"},
				providers: []string{"tls", "aws", "null"},
			},
		},
		{
			name: "sort module items",
			path: "full-example",
			sort: &SortBy{Name: true, Required: false},
			expected: expected{
				inputs:    []string{"A", "B", "C", "D", "E", "F"},
				required:  []string{"A", "F"},
				optional:  []string{"B", "C", "D", "E"},
				outputs:   []string{"A", "B", "C"},
				providers: []string{"aws", "null", "tls"},
			},
		},
		{
			name: "sort module items",
			path: "full-example",
			sort: &SortBy{Name: false, Required: true},
			expected: expected{
				inputs:    []string{"D", "B", "E", "A", "C", "F"},
				required:  []string{"A", "F"},
				optional:  []string{"D", "B", "E", "C"},
				outputs:   []string{"C", "A", "B"},
				providers: []string{"tls", "aws", "null"},
			},
		},
		{
			name: "sort module items",
			path: "full-example",
			sort: &SortBy{Name: true, Required: true},
			expected: expected{
				inputs:    []string{"A", "F", "B", "C", "D", "E"},
				required:  []string{"A", "F"},
				optional:  []string{"B", "C", "D", "E"},
				outputs:   []string{"A", "B", "C"},
				providers: []string{"aws", "null", "tls"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			path := filepath.Join("testdata", tt.path)
			options, _ := NewOptions().With(&Options{
				Path:   path,
				SortBy: tt.sort,
			})
			tfmodule, _ := loadModule(path)
			module, err := loadModuleItems(tfmodule, options)

			assert.Nil(err)
			sortItems(module, tt.sort)

			for i, v := range module.Inputs {
				assert.Equal(tt.expected.inputs[i], v.Name)
			}
			for i, v := range module.RequiredInputs {
				assert.Equal(tt.expected.required[i], v.Name)
			}
			for i, v := range module.OptionalInputs {
				assert.Equal(tt.expected.optional[i], v.Name)
			}
			for i, v := range module.Outputs {
				assert.Equal(tt.expected.outputs[i], v.Name)
			}
			for i, v := range module.Providers {
				assert.Equal(tt.expected.providers[i], v.Name)
			}
		})
	}
}
