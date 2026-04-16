import 'dart:io';

import 'package:google_mlkit_text_recognition/google_mlkit_text_recognition.dart';

class OcrService {
  final _recognizer = TextRecognizer(script: TextRecognitionScript.latin);

  Future<String?> extractPlate(File imageFile) async {
    final inputImage = InputImage.fromFile(imageFile);
    final RecognizedText result = await _recognizer.processImage(inputImage);

    if (result.text.isEmpty) return null;

    final List<String> blocks = result.blocks
        .expand((block) => block.lines)
        .map(((line) => line.text.trim().toUpperCase()))
        .where((text) => text.isNotEmpty)
        .toList();

    final plateCandidate = _findPlatePattern(blocks);
    if (plateCandidate != null) return plateCandidate;

    final fallback = blocks.firstWhere(
      (text) => text.length >= 4 && text.length <= 12,
      orElse: () => blocks.first,
    );

    return fallback;
  }

  String? _findPlatePattern(List<String> blocks) {
    final plateRegex = RegExp(r'^([A-Z]{1,2})\s*(\d{1,4})\s*([A-Z]{1,3})$');

    for (final block in blocks) {
      final normalized = block
          .replaceAll(RegExp(r'[^A-Z0-9\s]'), '')
          .replaceAll(RegExp(r'\s+'), ' ')
          .trim();

      if (plateRegex.hasMatch(normalized)) {
        final match = plateRegex.firstMatch(normalized);
        return '${match?.group(1)} ${match?.group(2)} ${match?.group(3)}';
      }
    }

    return null;
  }

  void dispose() {
    _recognizer.close();
  }
}
