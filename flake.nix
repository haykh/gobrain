{
  description = "a terminal-based notes and tasks organizer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      packages.${system}.default = pkgs.buildGoModule rec {
        name = "gobrain";
        author = "haykh";

        src = pkgs.fetchFromGitHub {
          owner = author;
          repo = name;
          rev = "master";
          hash = "sha256-KvG/b8NdzEMuj+zBL1+E/ItSb+6zXVXlLa2jqiNnvJE=";
        };

        vendorHash = "sha256-AfdgJceYwgZB0lVpVRWYoGkbBQWWCV+FHx33Ru7lCUM=";

        meta = with pkgs.lib; {
          homepage = "https://github.com/${author}/gobrain";
          license = licenses.unlicense;
          maintainers = [ author ];
        };
      };

      devShells.${system}.default = pkgs.mkShell {
        packages = [
          pkgs.go
          pkgs.gopls
          pkgs.gotools
        ];

        shellHook = ''
          echo "Welcome to gobrain development shell!"
        '';
      };
    };
}
