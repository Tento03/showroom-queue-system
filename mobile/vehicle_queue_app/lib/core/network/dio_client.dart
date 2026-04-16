import 'package:dio/dio.dart';

class DioClient {
  DioClient._();

  static final Dio instance = _buildDio();

  static Dio _buildDio() {
    final dio = Dio(
      BaseOptions(
        baseUrl: 'http://10.17.11.48:8080/',
        connectTimeout: const Duration(seconds: 10),
        receiveTimeout: const Duration(seconds: 15),
        headers: {'Accept': 'application/json'},
      ),
    );

    // ── Interceptors ────────────────────────────────────────────────────────
    // Add logging in debug mode only
    assert(() {
      dio.interceptors.add(
        LogInterceptor(
          requestBody: true,
          responseBody: true,
          requestHeader: false,
        ),
      );
      return true;
    }());

    return dio;
  }
}
