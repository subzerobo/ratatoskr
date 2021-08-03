# Ratatoskr (WIP)

![Ratatoskr](assets/art.png?raw=true "Ratatoskr")
Illustration by [Sergey Arzamastsev] (https://dribbble.com/arzarz)

Ratatoskr is a highly scalable push notification microservice written in Golang. 

## Project naming origin
[Ratatoskr](https://en.wikipedia.org/wiki/Ratatoskrhttps://en.wikipedia.org/wiki/Ratatoskr) (Old Norse, generally considered to mean "drill-tooth" or "bore-tooth") is a squirrel who runs up and down the world tree Yggdrasil to carry messages between the eagle perched atop Yggdrasil, and the serpent Níðhöggr, who dwells beneath one of the three roots of the tree.
Around it exists all else, including the Nine Worlds.

## Project Binaries
Ratatoskr is consisted of 4 binaries to make sure system is highly scalable in each part of project
- #### Yggdrasil 
  Back-Office API which is responsible for Accounting, Administration and Reporting of applications)
- #### Bifrost
  Public-facing API for Devices and Devices interaction
- #### Odin
  Responsible for managing notification and aggregation of results
- #### huggin
  Odin's workers responsible for interacting with Firebase API

## Installation

## Usage

###

## License
[MIT](https://choosealicense.com/licenses/mit/)