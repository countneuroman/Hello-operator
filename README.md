# Hello-operator

Простой оператор, реализованный с использованием [client-go](https://github.com/kubernetes/client-go) и [code-generator](https://github.com/kubernetes/code-generator)  

## Настройка окружения

1. go mod vendor для генерации папки vendor/ с зависимостями - требуется для использования code-generator
2. Запуск скрипта ./hack/update_codegen.sh для генерации кода (pkg/generated)