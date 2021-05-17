# Harborutils

This Go application helps propagating changes made in one Harbor registry to another. This changes can be new Oidc users, repositories permissions or robot accounts i.e. 

## License and Software Information
 
© adidas AG
 
adidas AG publishes this software and accompanied documentation (if any) subject to the terms of the MIT license with the aim of helping the community with our tools and libraries which we think can be also useful for other people. You will find a copy of the MIT license in the root folder of this package. All rights not explicitly granted to you under the MIT license remain the sole and exclusive property of adidas AG.
 
NOTICE: The software has been designed solely for the purpose of propagating changes from one harbor registry to another. The software is NOT designed, tested or verified for productive use whatsoever, nor or for any use related to high risk environments, such as health care, highly or fully autonomous driving, power plants, or other critical infrastructures or services.
 
If you want to contact adidas regarding the software, you can mail us at _software.engineering@adidas.com_.
 
For further information open the [adidas terms and conditions](https://github.com/adidas/adidas-contribution-guidelines/wiki/Terms-and-conditions) page.

Disclaimer
----------

adidas is not responsible for the usage of this software for different purposes that the ones described in the use cases.

Usage
-----

```
Usage:
  harborutils [command]

Available Commands:
  completion       generate the autocompletion script for the specified shell
  deleteGroups     Delete groups from Harbor
  fixEmptyEmails   Fix empty emails in database
  getGroups        Get groups from Harbor
  getProjects      Get projects from Harbor
  help             Help about any command
  importLdapGroups Propagate groups from primary harbor to secondary
  importLdapUsers  Propagate users from primary harbor to secondary
  syncGrants       Propagate grants from primary harbor to secondary
  syncLabels       Propagate project labels from primary harbor to secondary
  syncRobotAccount Propagate robot account from primary harbor to secundary
  syncUsersDb      Sync useres between harbor primarty and harbor secundary

Flags:
  -v, --apiVersion string   APIVersion (ie v2.0)
  -s, --harbor string       Harbor Server address
  -h, --help                help for harborutils
  -p, --password string     Password
  -u, --user string         Username Harbor

Use "harborutils [command] --help" for more information about a command.```
  

Releases
--------

* 1.0.0 - First version