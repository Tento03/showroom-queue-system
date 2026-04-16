import 'dart:io';

import 'package:vehicle_queue_app/core/network/dio_client.dart';
import 'package:vehicle_queue_app/features/queue/services/upload_service.dart';

class QueueService {
  final _dio = DioClient.instance;
  final _uploadService = UploadService();

  Future<String> createQueue({
    required String plate,
    required File image,
    required String ownerName,
    required String ownerPhone,
  }) async {
    final imageUrl = await _uploadService.uploadImage(image);

    final response = await _dio.post(
      '/queue',
      data: {
        'vehicle_plate': plate,
        'vehicle_image_url': imageUrl,
        'owner_name': ownerName,
        'owner_phone': ownerPhone,
      },
    );

    final queueNumber = response.data['queue_number'] as String?;

    if (queueNumber == null || queueNumber.isEmpty) {
      throw Exception('Queue created but queue_number missing in response');
    }
    return queueNumber;
  }
}
