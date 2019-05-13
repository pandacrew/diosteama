with import <nixpkgs> { };
with dockerTools;

let
diosteama = rustPlatform.buildRustPackage rec {
  name = "diosteama-${version}";
  version = "0.1";
  src = ./.;
  nativeBuildInputs = [ pkgconfig ];
  buildInputs = [ openssl ];
  cargoSha256 = "0shxn6fd3y9ma165d0i58l6v6mkjmr07x88gkfrv34c310c6y5px";
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
