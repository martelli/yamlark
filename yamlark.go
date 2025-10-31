package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	yaml "github.com/goccy/go-yaml"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"go.starlark.net/starlarkstruct"
	"github.com/BurntSushi/toml"
)

func interfaceToStarlarkValue(input interface{}) (starlark.Value, error) {
	if input == nil {
		return starlark.None, nil
	}

	switch v := input.(type) {
	case string:
		return starlark.String(v), nil
	case bool:
		return starlark.Bool(v), nil
	case int:
		return starlark.MakeInt(v), nil
	case int32:
		return starlark.MakeInt(int(v)), nil
	case int64:
		return starlark.MakeInt(int(v)), nil
	case uint64:
		return starlark.MakeInt(int(v)), nil
	case float32:
		return starlark.Float(v), nil
	case float64:
		return starlark.Float(v), nil

	case map[string]interface{}:
		dict := starlark.NewDict(len(v))
		for k, val := range v {
			starlarkKey := starlark.String(k)
			starlarkVal, err := interfaceToStarlarkValue(val)
			if err != nil {
				return nil, err
			}
			if err := dict.SetKey(starlarkKey, starlarkVal); err != nil {
				return nil, err
			}
		}
		return dict, nil

	case []interface{}:
		starlarkList := make([]starlark.Value, len(v))
		for i, val := range v {
			starlarkVal, err := interfaceToStarlarkValue(val)
			if err != nil {
				return nil, err
			}
			starlarkList[i] = starlarkVal
		}
		return starlark.NewList(starlarkList), nil

	default:
		return nil, fmt.Errorf("unsupported Go type for Starlark conversion: %s", reflect.TypeOf(input))
	}
}

func starlarkValueToInterface(value starlark.Value) (interface{}, error) {
	if value == nil || value == starlark.None {
		return nil, nil
	}

	switch v := value.(type) {
	case starlark.String:
		return v.GoString(), nil
	case starlark.Bool:
		return bool(v), nil
	case starlark.Int:
		i, ok := v.Int64()
		if !ok {
			return nil, fmt.Errorf("starlark int overflow")
		}
		return i, nil

	case *starlark.List:
		list := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			elem, err := starlarkValueToInterface(v.Index(i))
			if err != nil {
				return nil, err
			}
			list[i] = elem
		}
		return list, nil

	case starlark.Tuple:
		list := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			elem, err := starlarkValueToInterface(v.Index(i))
			if err != nil {
				return nil, err
			}
			list[i] = elem
		}
		return list, nil

	case *starlark.Dict:
		goMap := make(map[string]interface{})
		for _, item := range v.Items() {
			key := item[0]
			val := item[1]

			sKey, ok := key.(starlark.String)
			if !ok {
				return nil, fmt.Errorf("dictionary key is not a string: %s", key.Type())
			}

			goVal, err := starlarkValueToInterface(val)
			if err != nil {
				return nil, err
			}

			goMap[sKey.GoString()] = goVal
		}
		return goMap, nil

	default:
		return nil, fmt.Errorf("unsupported starlark type: %s", v.Type())
	}
}

func starlarkFileRead(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string

	if err := starlark.UnpackArgs("read", args, kwargs, "path", &path); err != nil {
		return starlark.None, fmt.Errorf("file.read: failed to unpack: %w", err)
	}

	baseDir, err := os.Getwd()
	if err != nil {
		return starlark.None, fmt.Errorf("file.read: could not get current directory: %w", err)
	}

	fullPath := filepath.Join(baseDir, path)
	cleanPath := filepath.Clean(fullPath)

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return starlark.None, fmt.Errorf("file.read failed for path %s: %w", path, err)
	}

	return starlark.String(data), nil
}

func starlarkFileWrite(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string
	var data string

	if err := starlark.UnpackArgs("write", args, kwargs, "path", &path, "data", &data); err != nil {
		return starlark.None, fmt.Errorf("file.write: failed to unpack: %w", err)
	}

	baseDir, err := os.Getwd()
	if err != nil {
		return starlark.None, fmt.Errorf("file.write: could not get current directory: %w", err)
	}

	fullPath := filepath.Join(baseDir, path)
	cleanPath := filepath.Clean(fullPath)

	err = os.WriteFile(cleanPath, []byte(data), 0644)
	if err != nil {
		return starlark.None, fmt.Errorf("file.write failed for path %s: %w", path, err)
	}

	return starlark.None, nil
}

func starlarkYamlDumps(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var data *starlark.Dict

	if err := starlark.UnpackArgs("dumps", args, kwargs, "data", &data); err != nil {
		return starlark.None, fmt.Errorf("yaml.dumps failed to unpack: %w", err)
	}

	bogus, err := starlarkValueToInterface(data)
	if err != nil {
		return starlark.None, fmt.Errorf("yaml.dumps failed to convert: %w", err)
	}

	bytes, err := yaml.Marshal(bogus)
	if err != nil {
		return starlark.None, err
	}

	return starlark.String(string(bytes)), nil
}

func starlarkYamlRead(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string
	var y interface{}

	if err := starlark.UnpackArgs("read", args, kwargs, "path", &path); err != nil {
		return starlark.None, fmt.Errorf("yaml.read failed to unpack: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("yaml.read failed to read path %q: %w", path, err)
	}

	err = yaml.Unmarshal(data, &y)
	if err != nil {
		return nil, fmt.Errorf("yaml.read failed to unmarshal path %q: %w", path, err)
	}

	dict, err := interfaceToStarlarkValue(y)
	if err != nil {
		return nil, fmt.Errorf("yaml.read failed to convert path %q: %w", path, err)
	}

	return dict, nil
}

func starlarkTomlRead(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string
	var tomlObj interface{}

	if err := starlark.UnpackArgs("read", args, kwargs, "path", &path); err != nil {
		return starlark.None, fmt.Errorf("toml.read: failed to unpack: %w", err)
	}

	baseDir, err := os.Getwd()
	if err != nil {
		return starlark.None, fmt.Errorf("toml.read: could not get current directory: %w", err)
	}

	fullPath := filepath.Join(baseDir, path)
	cleanPath := filepath.Clean(fullPath)

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return starlark.None, fmt.Errorf("toml.read failed for path %s: %w", path, err)
	}

	_, err = toml.Decode(string(data), &tomlObj)
	if err != nil {
		return starlark.None, fmt.Errorf("toml.read failed to parse path %s: %w", path, err)
	}

	dict, err := interfaceToStarlarkValue(tomlObj)
	if err != nil {
		return nil, fmt.Errorf("toml.read failed to convert path %q: %w", path, err)
	}

	return dict, nil
}

var FileModule = &starlarkstruct.Module{
	Name: "file",
	Members: starlark.StringDict{
		"read":  starlark.NewBuiltin("file.read", starlarkFileRead),
		"write": starlark.NewBuiltin("file.write", starlarkFileWrite),
	},
}

var YamlModule = &starlarkstruct.Module{
	Name: "yaml",
	Members: starlark.StringDict{
		"read": starlark.NewBuiltin("yaml.read", starlarkYamlRead),
		"dumps": starlark.NewBuiltin("yaml.dumps", starlarkYamlDumps),
	},
}

var TomlModule = &starlarkstruct.Module{
	Name: "toml",
	Members: starlark.StringDict{
		"read": starlark.NewBuiltin("toml.read", starlarkTomlRead),
	},
}

func getBuiltins() starlark.StringDict {
    return starlark.StringDict{
		"file":  FileModule,
		"yaml":  YamlModule,
		"toml":  TomlModule,
	}
}

func starlarkLoad(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	data, err := os.ReadFile(module)
	if err != nil {
		return nil, fmt.Errorf("load failed to read module %q: %w", module, err)
	}
	return starlark.ExecFileOptions(syntax.LegacyFileOptions(), thread, module, data, getBuiltins())
}

func executeStarlarkScript(filename string) error {
	thread := &starlark.Thread{Name: "main"}
	thread.Load = starlarkLoad
	_, err := starlark.ExecFileOptions(syntax.LegacyFileOptions(), thread, filename, nil, getBuiltins())

	if err != nil {
		return fmt.Errorf("starlark execution failed: %w", err)
	}
	return nil
}

func main() {

	if len(os.Args) != 2 {
		log.Fatal("\n--- Execution Error ---\nusage: lark <script>\n")
		return
	}

	if err := executeStarlarkScript(os.Args[1]); err != nil {
		log.Fatal("\n--- Execution Error ---\n%v\n", err)
	}
}
