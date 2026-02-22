{
  description = "Punchcard Go program";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
  };

  outputs = { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        rec {
          punchcard = pkgs.buildGoModule {
            pname = "punchcard";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-Ei1Rej/gkc5nuMgH2aR6qxqfEYWiADoS2s6b4sefGhk=";
            subPackages = [ "." ];

            postInstall = ''
              if [ -e "$out/bin/go-punchcard" ]; then
                mv "$out/bin/go-punchcard" "$out/bin/punchcard"
              fi
            '';
          };

          default = punchcard;
        });

      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            packages = [ pkgs.go ];
          };
        });
    };
}
