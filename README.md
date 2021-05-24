# lynx-example-integration
Example integration for exporting data from the Lynx platform using the [go-lynx](https://github.com/IoTOpen/go-lynx, "go-lynx") library.

This integration exports data from functions created by [Lynx-example-manager](https://github.com/IoTOpen/lynx-example-manager "Lynx-example-manager")
and will not work without the proper devices and functions that it sets up. Also note that import and export should be done simultaneously inorder to see any sort of result.

### Configure
Create a new configuration file based on the [lynx-integration.example.yml](lynx-integration.example.yml) and name it `lynx-integration.yml` and place it in the root folder of the project.