# Using vale.sh

Vale is a natural language linter that helps writers improve their prose. Monty uses it to standardize markdown files.

## Installing vale
Refer to: https://vale.sh/docs/vale-cli/installation/

The vale configuration file used for Monty is located at `.vale.ini`

After installing, run `vale sync` to fetch the styles specified in the configuration file. This only needs to be run once.

## Usage

Single file:
`vale path/to/file.md`

Multiple files:
`vale path/to/directory/`

## Docker

If you don't want to install vale, you can run it in a docker container by mounting the Monty source after installing the packages [`.vale.ini`](./.vale.ini):

```shell
docker run -v $(pwd):/monty jdkato/vale:v2.22.0 --config /monty/.vale.ini /monty/<file>
```
