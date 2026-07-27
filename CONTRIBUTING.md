# Adding a new resource

1. Confirm the resource has a list endpoint in its Magic Modules YAML (`base_url` with no `{{name}}`)
2. Run the generator:
   ```bash
   go run cmd/scaffold-list/main.go \
     --magic-modules ~/path/to/magic-modules \
     --resource <service>/<Resource>.yaml \
     --output /tmp
   ```
3. Copy the output into the provider and verify it builds:
   ```bash
   cp /tmp/data_source_<resource>_list.go \
     ~/path/to/terraform-provider-google/google/services/<service>/
   go build ./google/services/<service>/
   ```
4. Write an acceptance test modeled on `data_source_dns_managed_zone_list_test.go`
5. If the flattener needs type coercion (e.g. numeric string ID → int), add it to the template
