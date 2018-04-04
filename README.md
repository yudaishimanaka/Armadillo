<div align="center">
<img src="https://github.com/yudaishimanaka/Armadillo/blob/master/images/armadillo.png" alt="armadillo" width="128" height="128">
</div>

# Armadillo
[![Build Status](https://travis-ci.org/yudaishimanaka/Armadillo.svg?branch=master)](https://travis-ci.org/yudaishimanaka/Armadillo)

Simple password management CLI tool  
It enables you to manage passwords from the terminal, so you can manage passwords for each service.  

**It only supports password management to the last.**

## Install
1. Clone repository  
`~$ git clone github.com/yudaishimanaka/Armadillo.git`
2. Move binary  
`~$ sudo mv armadillo /usr/bin`
3. Grant execution authority  
`~$ sudo chmod +x /usr/bin/armadillo`  

## Usage
1. Initialize(First time only once)  
`~$ armadillo init`
2. See the help command  
`~$ armadillo help`  

## Comands
|Command|Detail|
|:--:|:--:|
|`init`|Initialization|
|`create`|Save information for the service|
|`delete`|Delete information for the service|
|`update`|Update information for the service|
|`show`|View service information (Service Name, Email or User ID, Password)|

## License
The MIT License (MIT) -see `LICENSE` for more details
