with import <nixpkgs> { };
with dockerTools;

let
diosteama = rustPlatform.buildRustPackage rec {
  name = "diosteama-${version}";
  version = "0.1";
  src = ./.;
  nativeBuildInputs = [ pkgconfig ];
  buildInputs = [ openssl ];
  cargoSha256 = "11iwd9fr7zyqgm5kainqighync90yp2lrsh1xal9yhriy3f6i2g5";
};

in
buildImage {
name = "diosteama";
  tag = "latest";
diskSize = 4096;
contents = [
  cacert
  diosteama
];
  config = {
    Cmd = [ "/bin/diosteama" ];
    Env = [
      "SSL_CERT_FILE=/etc/ssl/certs/ca-bundle.crt"
    ];
  };
}
