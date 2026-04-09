{
  description = "a terminal-based notes and tasks organizer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { stdenv.hostPlatform.system = system; };

        pname = "gobrain";
        version = "1.4.0";
        author = "haykh";

        src = pkgs.fetchFromGitHub {
          inherit pname version;
          owner = author;
          repo = pname;
          rev = "master";
          hash = "sha256-KGGgXZgy6Un/YW22Kbl39PMGyY2j7PmxhFPjCY24iyY=";
        };
      in
      {
        packages.default = pkgs.buildGoModule {
          inherit pname version src;

          vendorHash = "sha256-AfdgJceYwgZB0lVpVRWYoGkbBQWWCV+FHx33Ru7lCUM=";

          meta = with pkgs.lib; {
            homepage = "https://github.com/${author}/gobrain";
            license = licenses.unlicense;
            maintainers = [ author ];
          };
        };

        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.gopls
            pkgs.gotools
          ];

          shellHook = ''
            echo "Welcome to gobrain development shell!"
          '';
        };
      }
    );
}
