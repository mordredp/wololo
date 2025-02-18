module github.com/mordredp/wololo

go 1.20

require (
	github.com/go-ldap/ldap v3.0.3+incompatible // indirect
	github.com/ilyakaznacheev/cleanenv v1.4.2
)

//replace github.com/mordredp/auth => ../auth

require github.com/mordredp/auth v0.0.0-20230723154909-747f874e29c4

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/go-chi/chi/v5 v5.0.8
	github.com/joho/godotenv v1.5.1 // indirect
	gopkg.in/asn1-ber.v1 v1.0.0-20181015200546-f715ec2f112d // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
