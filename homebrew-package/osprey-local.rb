class OspreyLocal < Formula
  desc "A modern, systems programming language for building networked applications"
  homepage "https://ospreylang.com"
  version "0.2.0"
  license "MIT"

  if Hardware::CPU.arm?
    url "file://#{__dir__}/release/osprey-darwin-arm64.tar.gz"
    sha256 "0f7db9c85d7cf89240ffe8a09238f9d3cae9b0a6644260d9b2a1f20568020f41"
  else
    url "file://#{__dir__}/release/osprey-darwin-amd64.tar.gz"
    sha256 "c59c5259383b3eed3f1a8d80f6671c43ec7d90dbd78a4e9ad631284cde351924"
  end

  depends_on "llvm"

  def install
    bin.install "osprey"
    lib.install "libfiber_runtime.a"
    lib.install "libhttp_runtime.a"
    lib.install "libwebsocket_runtime.a"
    lib.install "libsystem_runtime.a"
  end

  test do
    # Create a simple Osprey program
    (testpath/"hello.osp").write <<~EOS
      fun main() -> Void {
        print("Hello, Homebrew!")
      }
    EOS

    # Test compilation
    system bin/"osprey", "compile", testpath/"hello.osp"
    
    # Check that osprey binary exists and is executable
    assert_predicate bin/"osprey", :exist?
    assert_predicate bin/"osprey", :executable?
  end
end 