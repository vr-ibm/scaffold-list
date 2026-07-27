# scaffold-list

A code generator that produces Terraform "list" data sources for the [terraform-provider-google](https://github.com/hashicorp/terraform-provider-google) from Magic Modules resource YAML definitions.

## What it does

Given a Magic Modules resource YAML (e.g. `dns/ManagedZone.yaml`), scaffold-list generates a complete, ready-to-compile Go file containing:

- **Schema** — a `TypeList` attribute (`managed_zones`, `instances`, etc.) with a nested `Elem` schema derived from the resource's properties, plus synthetic `id` and `project` fields
- **Paginated read function** — calls the resource's list endpoint via `transport_tpg.SendRequest`, follows `nextPageToken` until all pages are consumed
- **Flattener** — maps REST response keys (camelCase) to Terraform schema keys (snake_case), with type conversions where needed (e.g. string → int64 for numeric IDs)
- **Registry entry** — registers the data source via `registry.Schema` so the provider picks it up automatically

The generated file builds cleanly against terraform-provider-google and passes acceptance tests without manual edits.

## Usage

```bash
go run cmd/scaffold-list/main.go \
  --magic-modules ~/path/to/magic-modules \
  --resource <service>/<Resource>.yaml \
  --output /tmp
```

## Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--magic-modules` | yes | Path to a local clone of [magic-modules](https://github.com/GoogleCloudPlatform/magic-modules) |
| `--resource` | yes | Relative path to the resource YAML within `mmv1/products/` (e.g. `dns/ManagedZone.yaml`) |
| `--output` | yes | Directory where the generated `.go` file will be written |

## Example

```bash
# Generate the dns managed zone list data source
go run cmd/scaffold-list/main.go \
  --magic-modules ~/Documents/github/magic-modules \
  --resource dns/ManagedZone.yaml \
  --output /tmp

# Copy into the provider and build
cp /tmp/data_source_managed_zone_list.go \
  ~/Documents/github/terraform-provider-google/google/services/dns/

cd ~/Documents/github/terraform-provider-google
go build ./google/services/dns/
```

## Generated output

For `dns/ManagedZone.yaml`, the tool produces `data_source_managed_zone_list.go` which registers `google_dns_managed_zone_list` and exposes:

| Attribute | Type | Description |
|---|---|---|
| `id` | string | `projects/{project}/managedZones` |
| `project` | string | The GCP project (from provider config or explicit) |
| `managed_zones` | list | All zones in the project |
| `managed_zones.*.id` | string | `projects/{project}/managedZones/{name}` |
| `managed_zones.*.project` | string | GCP project |
| `managed_zones.*.name` | string | Zone name |
| `managed_zones.*.dns_name` | string | DNS name (e.g. `example.com.`) |
| `managed_zones.*.description` | string | Zone description |
| `managed_zones.*.visibility` | string | `public` or `private` |
| `managed_zones.*.managed_zone_id` | int | Numeric zone ID |

## How it works

1. Parses the Magic Modules resource YAML to extract:
   - API base URL and list endpoint
   - Resource properties (names, types, nesting)
   - Response body field name (e.g. `managedZones`)

2. Generates the Terraform schema from resource properties, converting MM types to `schema.Type*`

3. Emits a read function with pagination loop using `nextPageToken`

4. Emits a flattener that maps each property's API name to its snake_case schema name, with type coercion for fields where the REST JSON type doesn't match the schema type (e.g. numeric string IDs → int64)

5. Emits a `registry.Schema` `init` block for automatic provider registration

## Requirements

- Go 1.21+
- A local clone of <https://github.com/GoogleCloudPlatform/magic-modules>
- A local clone of <https://github.com/hashicorp/terraform-provider-google> (for building/testing output)

## Testing generated output

Write an acceptance test following the provider's conventions:

```go
func TestAccDataSourceDnsManagedZoneList_basic(t *testing.T) {
    // Create resources, then read via data source
    // Assert managed_zones.# > 0
    // Assert managed_zones.0.name is set
}
```

Run with:

```bash
TF_ACC=1 go test ./google/services/dns/ \
  -run TestAccDataSourceDnsManagedZoneList_basic \
  -v -timeout 300s
```
