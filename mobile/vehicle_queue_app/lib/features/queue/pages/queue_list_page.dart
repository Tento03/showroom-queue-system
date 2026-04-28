import 'package:flutter/material.dart';
import 'package:vehicle_queue_app/core/network/dio_client.dart';

class QueueListPage extends StatefulWidget {
  const QueueListPage({super.key});

  @override
  State<QueueListPage> createState() => _QueueListPageState();
}

class _QueueListPageState extends State<QueueListPage> {
  final _dio = DioClient.instance;
  List<Map<String, dynamic>> _queues = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _fetchQueues();
  }

  Future<void> _fetchQueues() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });
    try {
      final response = await _dio.get('/queue');
      final data = response.data;
      List<dynamic> list = [];
      if (data is List) {
        list = data;
      } else if (data is Map && data['data'] is List) {
        list = data['data'] as List;
      }
      setState(() {
        _queues = list.map((e) => Map<String, dynamic>.from(e as Map)).toList();
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F7FA),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        centerTitle: false,
        title: const Text(
          'Daftar Antrian',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600),
        ),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(1),
          child: Divider(height: 1, color: Colors.grey.shade200),
        ),
        actions: [
          IconButton(
            onPressed: _fetchQueues,
            icon: const Icon(Icons.refresh_rounded),
            tooltip: 'Refresh',
          ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            CircularProgressIndicator(color: Color(0xFF1A73E8)),
            SizedBox(height: 16),
            Text(
              'Memuat daftar antrian...',
              style: TextStyle(color: Colors.grey, fontSize: 14),
            ),
          ],
        ),
      );
    }

    if (_error != null) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  color: Colors.red.shade50,
                  shape: BoxShape.circle,
                ),
                child: Icon(
                  Icons.wifi_off_rounded,
                  size: 40,
                  color: Colors.red.shade400,
                ),
              ),
              const SizedBox(height: 16),
              const Text(
                'Gagal memuat data',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
              ),
              const SizedBox(height: 8),
              Text(
                'Periksa koneksi internet atau coba lagi.',
                style: TextStyle(fontSize: 13, color: Colors.grey.shade600),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 20),
              FilledButton.icon(
                onPressed: _fetchQueues,
                icon: const Icon(Icons.refresh_rounded, size: 18),
                label: const Text('Coba Lagi'),
                style: FilledButton.styleFrom(
                  backgroundColor: const Color(0xFF1A73E8),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                ),
              ),
            ],
          ),
        ),
      );
    }

    if (_queues.isEmpty) {
      return Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              padding: const EdgeInsets.all(20),
              decoration: const BoxDecoration(
                color: Color(0xFFE8F0FE),
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.format_list_bulleted_rounded,
                size: 40,
                color: Color(0xFF1A73E8),
              ),
            ),
            const SizedBox(height: 16),
            const Text(
              'Belum ada antrian',
              style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
            ),
            const SizedBox(height: 8),
            Text(
              'Antrian yang dibuat akan muncul di sini.',
              style: TextStyle(fontSize: 13, color: Colors.grey.shade500),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _fetchQueues,
      color: const Color(0xFF1A73E8),
      child: ListView.separated(
        padding: const EdgeInsets.all(16),
        itemCount: _queues.length,
        separatorBuilder: (_, __) => const SizedBox(height: 10),
        itemBuilder: (context, index) {
          final q = _queues[index];
          return _QueueCard(data: q);
        },
      ),
    );
  }
}

// ─── Queue Card ──────────────────────────────────────────────────────────────

class _QueueCard extends StatelessWidget {
  const _QueueCard({required this.data});
  final Map<String, dynamic> data;

  @override
  Widget build(BuildContext context) {
    final queueNumber = data['queue_number']?.toString() ?? '-';
    final plate = data['vehicle_plate']?.toString() ?? '-';
    final ownerName = data['owner_name']?.toString() ?? '-';
    final ownerPhone = data['owner_phone']?.toString() ?? '-';
    final status = data['status']?.toString() ?? 'WAITING';
    final imageUrl = data['vehicle_image_url']?.toString();

    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // ── Queue Number Badge ────────────────────────────────────────────
            Container(
              width: 56,
              height: 56,
              decoration: BoxDecoration(
                gradient: const LinearGradient(
                  colors: [Color(0xFF1A73E8), Color(0xFF0D47A1)],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(14),
              ),
              alignment: Alignment.center,
              child: Text(
                queueNumber,
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.w800,
                  fontSize: 14,
                  letterSpacing: 1,
                ),
              ),
            ),
            const SizedBox(width: 14),

            // ── Info ─────────────────────────────────────────────────────────
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Expanded(
                        child: Text(
                          plate,
                          style: const TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w700,
                            letterSpacing: 1.5,
                          ),
                        ),
                      ),
                      _StatusBadge(status: status),
                    ],
                  ),
                  const SizedBox(height: 6),
                  _InfoRow(
                    icon: Icons.person_outline,
                    value: ownerName,
                  ),
                  const SizedBox(height: 4),
                  _InfoRow(
                    icon: Icons.phone_outlined,
                    value: ownerPhone,
                  ),
                ],
              ),
            ),

            // ── Vehicle Thumbnail ─────────────────────────────────────────────
            if (imageUrl != null && imageUrl.isNotEmpty) ...[
              const SizedBox(width: 12),
              ClipRRect(
                borderRadius: BorderRadius.circular(10),
                child: Image.network(
                  imageUrl,
                  width: 56,
                  height: 56,
                  fit: BoxFit.cover,
                  errorBuilder: (_, __, ___) => Container(
                    width: 56,
                    height: 56,
                    decoration: BoxDecoration(
                      color: Colors.grey.shade100,
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: Icon(
                      Icons.directions_car_outlined,
                      color: Colors.grey.shade400,
                    ),
                  ),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  const _InfoRow({required this.icon, required this.value});
  final IconData icon;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, size: 13, color: Colors.grey.shade400),
        const SizedBox(width: 5),
        Expanded(
          child: Text(
            value,
            style: TextStyle(fontSize: 13, color: Colors.grey.shade600),
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({required this.status});
  final String status;

  @override
  Widget build(BuildContext context) {
    final upper = status.toUpperCase();

    Color bg;
    Color fg;
    IconData icon;
    String label;

    switch (upper) {
      case 'WAITING':
        bg = const Color(0xFFFFF8E1);
        fg = const Color(0xFFF9A825);
        icon = Icons.hourglass_top_rounded;
        label = 'Menunggu';
        break;
      case 'IN_SERVICE':
      case 'IN SERVICE':
        bg = const Color(0xFFE3F2FD);
        fg = const Color(0xFF1565C0);
        icon = Icons.build_circle_outlined;
        label = 'Servis';
        break;
      case 'DONE':
      case 'COMPLETED':
        bg = const Color(0xFFE6F4EA);
        fg = const Color(0xFF2E7D32);
        icon = Icons.check_circle_outline;
        label = 'Selesai';
        break;
      case 'CANCELLED':
        bg = const Color(0xFFFCE4EC);
        fg = const Color(0xFFC62828);
        icon = Icons.cancel_outlined;
        label = 'Batal';
        break;
      default:
        bg = Colors.grey.shade100;
        fg = Colors.grey.shade600;
        icon = Icons.info_outline;
        label = status;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: bg,
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 12, color: fg),
          const SizedBox(width: 4),
          Text(
            label,
            style: TextStyle(
              fontSize: 11,
              fontWeight: FontWeight.w600,
              color: fg,
            ),
          ),
        ],
      ),
    );
  }
}
