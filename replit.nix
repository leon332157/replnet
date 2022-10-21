{ pkgs }: {
	deps = [
    pkgs.netcat-gnu
    pkgs.go_1_18
	pkgs.python310
    pkgs.python310Packages.autopep8
	];
}