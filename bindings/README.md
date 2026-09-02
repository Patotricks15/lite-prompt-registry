# SDKs

Python packages use PyO3 and maturin; TypeScript packages use napi-rs; Go packages link the core through a narrow cgo ABI. Each SDK exposes prompt creation and immutable version creation without leaking storage details.