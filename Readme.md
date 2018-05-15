# ![Terraform provider for NCloud](./logo.jpg)

[Terraform](https://www.terraform.io/) provider for
[Naver Cloud Platform](https://www.ncloud.com/) (also known as NCloud).

## Installation

Simply [download the release for your target platform](./releases), and
place it into one of the following locations:

  1. On Windows, in the sub-path terraform.d/plugins beneath your 
     user's "Application Data" directory.
  2. On all other systems, in the sub-path .terraform.d/plugins in your
     user's home directory.

See https://www.terraform.io/docs/configuration/providers.html#third-party-plugins for more details.

## Usage

> ./main.tf

```hcl
provider "ncloud" {
  access_key = "$NCLOUD_ACCESS_KEY"
  secret_key = "$NCLOUD_SECRET_KEY"
  region     = "$NCLOUD_DEFAULT_REGION"
}

resource "ncloud_instance" "gameserver" {
  # coming zoon 
}
```

See [Services.md](./Services.md) for a detailed lists of available
regions, zones, servers and images.

## Development

### Requirements

- [Go](https://golang.org/) 1.8 or higher
- [dep](https://github.com/golang/dep) for dependency management
- [mmake](https://github.com/tj/mmake) to build this project
- [github-release](https://github.com/aktau/github-release) to make releases

### Setup

> Fork this project, then clone your fork.

```shell
git clone git@github.com:[username]/terraform-provider-ncloud.git
cd terraform-provider-ncloud
git remote add upstream git@github.com:Wizcorp/terraform-provider-ncloud.git
```

> Install dependencies, and build the project

```shell
dep ensure
mmake
```

See `mmake help` for more project-related commands.

> Update the list of code for server types and images

```shell
# *nix, macOS
export NCLOUD_ACCESS_KEY=""
export NCLOUD_SECRET_KEY=""

# Windows
set-item env:NCLOUD_ACCESS_KEY=""
set-item env:NCLOUD_SECRET_KEY=""

mmake generate-services
```

This command will update the `Service.md` file. This file contains a list
of relevant codes you will need to put in your Terraform configuration.

> Release

```shell
# *nix, macOS
export GITHUB_TOKEN="..."

# Windows
set-item env:GITHUB_TOKEN="..."

git tag v1.2.3
mmake release
```

See https://help.github.com/articles/creating-an-access-token-for-command-line-use to learn
how to create your token.

## Acknowledgements

Inspiration:

- [terraform-provider-linode](https://github.com/btobolaski/terraform-provider-linode)
- [terraform-provider-aws](https://github.com/terraform-providers/terraform-provider-aws)

## License

MIT. [See License](./License.md).