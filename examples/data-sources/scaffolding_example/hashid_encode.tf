data "hashid_encode" "example" {
  salt         = "dev"
  alphabet     = "abcdefghijklmnopqrstuvwxyz0123456789"
  min_length   = 30
  encode_value = "testvalue"

}

output "hashid" {
  value = data.hashid_encode.example.hash_id
}