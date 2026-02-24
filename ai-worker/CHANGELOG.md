# Changelog

All notable changes to the AI Manga Translator project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [10.1.0] - 2026-02-24

### Added

- `manga-ocr==0.1.13` pinned in `requirements.txt` for stability against upstream regressions

### Changed

- **MangaOCR initialized with `force_cpu=True`** — avoids a CUDA context conflict with `llama-cpp-python` on Windows when both are loaded in the same process (symptoms: segfault or `exit status 0xc0000005` after `✅ LLM loaded.`)
- **AI pipeline loading order**: LLM is now loaded first (before YOLO and MangaOCR) to establish the llama.cpp CUDA context before PyTorch initializes

### Fixed

- **Pipeline crash when spawned by Go worker on Windows**: Python subprocess fell back to `cp1252` encoding for piped stdout; emoji characters (`⏳`, `✅`) caused a fatal `UnicodeEncodeError`. Fixed on the Go side by setting `PYTHONIOENCODING=utf-8`, and on the Python side by ensuring no encoding-sensitive output happens before the env var is applied.
- **MangaOCR / transformers ViT segfault**: Third-party package installation (e.g. `IOPaint`, `diffusers`, `accelerate`) into the same venv could corrupt the shared GPU context used by manga-ocr's transformers model. `force_cpu=True` is the stable workaround.

## [10.0.0] - 2025-12-08

### Added

- **Intelligent Masked Inpainting**: Revolutionary text cleaning approach that preserves artwork and backgrounds
  - New `clean_box()` method using OpenCV masked inpainting instead of rectangle erasure
  - Fixed threshold-based text detection using `cv2.threshold` with configurable threshold value
  - Morphological dilation with `cv2.dilate` to expand text mask and cover anti-aliasing
  - TELEA inpainting algorithm (`cv2.inpaint`) for intelligent pixel filling
- New configuration parameters in `config/settings.py`:
  - `INPAINT_RADIUS` (default: 3) - Controls the radius for cv2.inpaint algorithm
  - `INPAINT_DILATE_ITERATIONS` (default: 1) - Number of dilation passes for text mask expansion
  - `INPAINT_DILATE_KERNEL_SIZE` (default: 3) - Kernel size for dilation operation (3x3)
  - `INPAINT_TEXT_THRESHOLD` (default: 180) - Threshold for dark text detection (0-255, lower = more aggressive)

### Changed

- **BREAKING**: Complete rewrite of `clean_box()` method in `services/typesetting.py`
  - Replaced simple rounded rectangle erasure with intelligent masked inpainting
  - Method now detects only dark text pixels instead of erasing entire bounding boxes
  - Improved handling of overlapping bounding boxes
- Updated pipeline version string to "V10 - Stable | Masked Inpainting"
- Enhanced README.md with V10 features and new configuration options

### Fixed

- Artifacts when bounding boxes overlap or merge in close proximity
- Background and artwork damage in regions between adjacent speech bubbles
- Text cleaning issues in complex bubble arrangements

## [9.0.0] - 2025

### Added

- Smart Box Merging algorithm to consolidate fragmented vertical text bubbles into single coherent blocks
- Custom "Anti-Thinking" prompt engineering to prevent LLM hallucinations
- Enhanced text sanitization to filter unsupported characters (emojis, complex symbols)
- Pixel-perfect text wrapping algorithm using pixel width measurement instead of character count

### Changed

- Improved batch processing with native ZIP support
- Format normalization to automatically convert all outputs to high-quality JPG
- Modular architecture with clear separation of concerns

### Fixed

- Text overflow and overlapping issues in narrow speech bubbles
- Font glitches caused by unsupported characters
- Internal monologues appearing in translated output

## [8.0.0] - 2025

### Added

- Initial public release
- YOLOv8-based text detection (fine-tuned on Manga109)
- MangaOCR integration for Japanese text recognition
- Qwen 2.5 7B (Abliterated) for uncensored translation
- Basic typesetting with rounded rectangle cleaning
- GPU acceleration via llama.cpp and CUDA
- Single image and ZIP batch processing

---

## Version History Summary

- **v10.1.0** - Windows subprocess fix, MangaOCR CPU mode, pipeline load order (Current)
- **v10.0.0** - Intelligent Masked Inpainting
- **v9.0.0** - Smart Box Merging & Anti-Thinking Prompts
- **v8.0.0** - Initial Public Release

[Unreleased]: https://github.com/P4ST4S/manga-translator/compare/v10.0.0...HEAD
