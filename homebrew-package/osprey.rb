class Osprey < Formula
  desc "Modern functional programming language designed for clarity, safety, and expressiveness"
  homepage "https://www.ospreylang.dev"
  url "https://github.com/melbournedeveloper/osprey/releases/download/v0.2.0/osprey-darwin-amd64.tar.gz"
  version "0.2.0"
  sha256 "ab9e08edc663ff82c22f2c9ec243122f6cb8bacebb48629b1c2237f48063fb40"
  
  depends_on "llvm"

  livecheck do
    url :stable
    regex(/^v?(\d+(?:\.\d+)+)$/i)
  end

  def install
    # Install pre-built binaries and ALL 4 runtime libraries
    bin.install "osprey"
    lib.install "libfiber_runtime.a"
    lib.install "libhttp_runtime.a"
    lib.install "libwebsocket_runtime.a"
    lib.install "libsystem_runtime.a"
  end

  test do
    # Test that the compiler can show help
    output = shell_output("#{bin}/osprey --help 2>&1", 0)
    assert_match "Osprey", output
    
    # Test that ALL 4 runtime libraries are installed
    assert_predicate lib/"libfiber_runtime.a", :exist?
    assert_predicate lib/"libhttp_runtime.a", :exist?
    assert_predicate lib/"libwebsocket_runtime.a", :exist?
    assert_predicate lib/"libsystem_runtime.a", :exist?
  end
end 