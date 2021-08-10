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
Interacts with harbor registry API

Usage:
  harborutils [command]

Available Commands:
  checkSha         Check image digest against Harbor
  deleteGroups     Delete groups from Harbor
  fixEmptyEmails   Fix empty emails in database
  getGroups        Get groups from Harbor
  getProjects      Get projects from Harbor
  getSha           Get image digest from Harbor
  help             Help about any command
  importLdapGroups Propagate groups from primary harbor to secondary
  importLdapUsers  Propagate users from primary harbor to secondary
  replicationTaks  returns the status of the last replications tasks (harbor stored the last 50)
  server           Run a server exposing some options of the cli
  syncGrants       Propagate grants from primary harbor to secondary
  syncLabels       Propagate project labels from primary harbor to secondary
  syncRegistries   Syncs objects created between two dates from harbor primary to harbor secundary
  syncRobotAccount Propagate robot account from primary harbor to secondary
  syncUsersDb      Sync users between harbor primary and harbor secondary

Flags:
  -v, --apiVersion string   APIVersion (ie v2.0) (default "v2.0")
  -s, --harbor string       Harbor Server address
  -h, --help                help for harborutils
  -p, --password string     Password
  -u, --user string         Username Harbor

Use "harborutils [command] --help" for more information about a command.
```

Update swagger documentation
----------------------------
```
$ cd server
$ swag  init -g root.go 
2021/06/15 20:13:47 Generate swagger docs....
2021/06/15 20:13:47 Generate general API Info, search dir:./
2021/06/15 20:13:47 Generating server.Token
2021/06/15 20:13:47 Generating server.APIError
2021/06/15 20:13:47 Generating server.ArtifactSha
2021/06/15 20:13:47 Generating server.ArtifactCheckSha
2021/06/15 20:13:47 create docs.go at docs/docs.go
2021/06/15 20:13:47 create swagger.json at docs/swagger.json
2021/06/15 20:13:47 create swagger.yaml at docs/swagger.yaml
```

Server Usage
------------
Api documentation in: http://localhost:8080/swagger/index.html

Examples written in [Httpie|https://httpie.io/]
### Get token

```
http -a "MyAzureUser:MyPasswordUser"  "http://localhost:8080/jwt"
```
### Get Image SHA

```
http "http://localhost:8080/artifact/sha"  image=="pea-cicd/test/debian:stable-20200607-slim" -a "MyAzureUser:MyPasswordUser"
http "http://localhost:8080/artifact/sha"  Token:MyToken image=="pea-cicd/test/debian:stable-20200607-slim"
```

### Check Image SHA

```
http "http://localhost:8080/artifact/check_sha"  image=="pea-cicd/test/debian:stable-20200607-slim" targetDigest==sha256:a1c2d5c775a3b7ebc7af29c77241819a86cd1222b1931d0712afdcd69c7dcbd5 -a "MyAzureUser:MyPasswordUser"
http "http://localhost:8080/artifact/check_sha"  Token:MyToken image=="pea-cicd/test/debian:stable-20200607-slim" targetDigest==sha256:a1c2d5c775a3b7ebc7af29c77241819a86cd1222b1931d0712afdcd69c7dcbd5
```


### Get Config

```
http "http://localhost:8080/config"
```

### Get Health

```
http "http://localhost:8080/health"
```


Releases
--------

* 1.0.0 - First version
* 1.1.0 - Add server
* 1.1.1 - fix oidctocket request
* 1.1.2 - Add error message when job can't run the replication in  replicationTaks command
