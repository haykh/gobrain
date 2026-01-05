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
        pname = "gobrain";
        version = "1.0.2";
        author = "haykh";

        src = pkgs.fetchFromGitHub {
          owner = author;
          repo = pname;
          rev = "v${version}";
          hash = "sha256-LP2q6O+2ApAC3VWjGKbviVJQ6uLwFOkAvdtFRHiSJ4c=";
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
