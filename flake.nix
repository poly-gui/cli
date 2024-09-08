{
  description = "Poly CLI";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs?tag=24.05";
  };

  outputs = { nixpkgs, ... }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          poly-cli = pkgs.buildGoModule {
            pname = "poly-cli";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-5USKR8JDAwMPOmef6rp93BQgBbdFz8xHt1+jLfTKp5U=";
          };

          default = poly-cli;
        });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.go
              pkgs.gotools
            ];
          };
        });
    };
}
