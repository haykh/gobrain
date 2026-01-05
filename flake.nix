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
        version = "1.0.1";
        author = "haykh";

        src = pkgs.fetchFromGitHub {
          owner = author;
          repo = pname;
          rev = "v${version}";
          hash = "sha256-Ye51PZ3jLUfaeD1iffrESFbXRzKPDJUEezVpqlk8EXs=";
        };

        vendorHash = "sha256-7mqgCcfy+VTOnCAPJmHWnVQL/7KAAzlVh0aHh/D0u4I=";

        meta = with pkgs.lib; {
          description = "go-based tool to do awesome stuff with notion";
          homepage = "https://github.com/${author}/nogo";
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
