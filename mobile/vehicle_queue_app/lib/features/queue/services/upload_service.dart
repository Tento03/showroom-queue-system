import 'dart:io';

import 'package:dio/dio.dart';
import 'package:vehicle_queue_app/core/network/dio_client.dart';

class UploadService {
  final _dio = DioClient.instance;

  Future<String> uploadImage(File imageFile) async {
    final formData = FormData.fromMap({
      'image': await MultipartFile.fromFile(
        imageFile.path,
        filename: 'vehicle_${DateTime.now().millisecondsSinceEpoch}.jpg',
      ),
    });

    final response = await _dio.post('/upload', data: formData);

    final imageUrl = response.data['image_url'] as String?;

    if (imageUrl == null || imageUrl.isEmpty) {
      throw Exception("Upload succeeded but image_url is missing in response");
    }

    return imageUrl;
  }
}
