{
  pkgs ? import <nixpkgs> { },
}:

let
  name = "gopop";
in
pkgs.mkShell ({
  name = "${name}-env";
  nativeBuildInputs = with pkgs; [
    go
    gopls
    gotools
    hugo
    wails
    prettier-plugin-go-template

    pkg-config
    makeWrapper
    gtk3
    webkitgtk_4_0

    nodePackages.nodejs
    vscode-langservers-extracted
    emmet-ls
    typescript-language-server
    taplo
    yaml-language-server
    markdown-oxide
    prettierd
    eslint_d
    mdformat
    svelte-language-server
  ];

  LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath [
    pkgs.stdenv.cc.cc
    pkgs.zlib
  ];
})
