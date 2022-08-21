package environment

import (
    "fmt"
    )

type Environment struct {
    values map[string]interface{}
}

func New() Environment {
    return Environment{values: make(map[string]interface{})}
}

func (env *Environment) Define(name string, value interface{}) {
    env.values[name] = value
}

func (env Environment) Get(name string) (interface{}, error) {
    val, ok := env.values[name]
    if !ok {
        return nil, fmt.Errorf("undefined variable '%s'", name)
    }
    return val, nil
}

func (env Environment) Assign(name string, value interface{}) error {
    if env.Exists(name) {
        env.values[name] = value
        return nil
    }
    return fmt.Errorf("undefined variable, '%s'", name)
}

func (env Environment) Exists(name string) bool {
    _, exists := env.values[name]
    return exists
}
