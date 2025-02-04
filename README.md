# Hello-operator

Простой оператор, реализованный с использованием [client-go](https://github.com/kubernetes/client-go) и [code-generator](https://github.com/kubernetes/code-generator)  

При применении нашего CRD создается job, который запускает pod, который выводит строку, указанную в нашем CRD в пармметре `message`

## Настройка окружения для разработки

* `go mod vendor` для генерации папки `vendor/` с зависимостями - требуется для использования code-generator
* Запуск скрипта `/hack/update_codegen.sh`  для генерации кода. В нашем случае он сгенерирован и распологается в  `/pkg/generated`  

## Запуск
1. Билдим контроллер: `go build -o hello-controller . `
2. Добавляем CRD в кластер `kubectl create -f crds/echo.yaml`
3. Запускаем контроллер `./hello-controller -kubeconfig=path-to-your-cluser-config.yaml`   
Запусается не в виде пода, а прямо на нашей локальной машине, удаленно подключаясь к кластеру.
4. Применяем пример нашего CRD `kubectl create -f crds/examples/echo.yaml`

## Полезные ссылки
[Официальный пример контроллера ](https://github.com/kubernetes/sample-controller)  
[Пример оператора](https://github.com/mmontes11/echoperator)  
[Заметки по kubernetes (Китайский)](https://github.com/huweihuang/kubernetes-notes)  
[SDK для создания операторов](https://github.com/kubernetes-sigs/kubebuilder)  
[Deep dive в генерацию на основе code-generat0r](https://www.redhat.com/en/blog/kubernetes-deep-dive-code-generation-customresources)