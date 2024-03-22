# Home Inventory
A small project to help us keep track of our crap. Used to keep track of containers in storage rooms.

> [!NOTE]
> Home Inventory is still in early alpha. Many planned features are not yet implemented, but you can still use the project to view the information.
> As of now, data is read-only. You need to edit the JSON files manually to use Home Inventory.

## Features
- Easy-to-use Web Interface
- Storage unit/locations
- Different types of locations

## Planned Features
- QR Code Reader
- QR Code Generator
- Search
- Create location/container
- Themes

## Self-hosting
It's very easy to start up your own Home Inventory server! As of now, Home Inventory only supports Linux.

Clone this repository as-is and run `go run .` in the project's root directory. No compilation necessary!
Home Inventory uses `/var/lib/HomeInventory` as its library directory. Its JSON files are stored here.

## Command line arguments
When starting Home Inventory, there are two arguments you can pass. They are both optional and can be used together.

- `-port [port]` - Specify a port to run Home Inventory (Default 8080)
- `-dev` - Start in dev mode (uses the working directory as the library)
