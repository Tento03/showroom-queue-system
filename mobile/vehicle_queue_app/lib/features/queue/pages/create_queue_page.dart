import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:image_picker/image_picker.dart';

import '../services/queue_service.dart';

/// UI + QueueService — tombol Scan tetap ada di UI tapi tidak terhubung ke OCR.
/// Submit sepenuhnya terhubung ke [QueueService].
class CreateQueuePage extends StatefulWidget {
  const CreateQueuePage({super.key});

  @override
  State<CreateQueuePage> createState() => _CreateQueuePageState();
}

class _CreateQueuePageState extends State<CreateQueuePage> {
  // ─── Controllers ──────────────────────────────────────────────────────────
  final _formKey = GlobalKey<FormState>();
  final _plateController = TextEditingController();
  final _ownerNameController = TextEditingController();
  final _ownerPhoneController = TextEditingController();

  // ─── Services ─────────────────────────────────────────────────────────────
  final _queueService = QueueService();

  // ─── State ────────────────────────────────────────────────────────────────
  File? _vehicleImage;
  bool _isLoading = false;

  // ─── Lifecycle ────────────────────────────────────────────────────────────
  @override
  void dispose() {
    _plateController.dispose();
    _ownerNameController.dispose();
    _ownerPhoneController.dispose();
    super.dispose();
  }

  // ─── Logic ────────────────────────────────────────────────────────────────

  Future<void> pickVehicleImage() async {
    final picked = await ImagePicker().pickImage(
      source: ImageSource.camera,
      imageQuality: 60,
    );
    if (picked != null) {
      setState(() => _vehicleImage = File(picked.path));
    }
  }

  /// Scan button membuka kamera tapi tidak memproses OCR.
  /// User tetap harus mengetik plat secara manual setelah foto diambil.
  Future<void> scanPlate() async {
    // Memberikan feedback visual bahwa tombol diklik dengan membuka picker (tanpa OCR)
    await ImagePicker().pickImage(source: ImageSource.camera, imageQuality: 50);

    // _showSnackBar(
    //   'Fitur baca plat otomatis sedang dikembangkan. Silakan ketik nomor plat secara manual.',
    //   isError: false,
    //   icon: Icons.info_outline,
    // );
  }

  Future<void> submit() async {
    if (!_formKey.currentState!.validate()) return;

    if (_vehicleImage == null) {
      _showSnackBar(
        'Harap ambil foto kendaraan terlebih dahulu.',
        isError: true,
        icon: Icons.camera_alt_outlined,
      );
      return;
    }

    setState(() => _isLoading = true);

    try {
      final queueNumber = await _queueService.createQueue(
        plate: _plateController.text.trim(),
        image: _vehicleImage!,
        ownerName: _ownerNameController.text.trim(),
        ownerPhone: _ownerPhoneController.text.trim(),
      );

      if (!mounted) return;
      _showQueueSuccessDialog(queueNumber);
      _resetForm();
    } catch (e) {
      _showSnackBar('Gagal submit: $e', isError: true);
    } finally {
      setState(() => _isLoading = false);
    }
  }

  void _resetForm() {
    _plateController.clear();
    _ownerNameController.clear();
    _ownerPhoneController.clear();
    setState(() => _vehicleImage = null);
  }

  void _showQueueSuccessDialog(String queueNumber) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (_) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        contentPadding: const EdgeInsets.fromLTRB(24, 28, 24, 20),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              padding: const EdgeInsets.all(16),
              decoration: const BoxDecoration(
                color: Color(0xFFE6F4EA),
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.check_circle_outline,
                color: Color(0xFF34A853),
                size: 40,
              ),
            ),
            const SizedBox(height: 16),
            const Text(
              'Kendaraan Ditambahkan!',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.w700),
            ),
            const SizedBox(height: 8),
            const Text(
              'Nomor Antrian',
              style: TextStyle(fontSize: 13, color: Colors.grey),
            ),
            const SizedBox(height: 6),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
              decoration: BoxDecoration(
                color: const Color(0xFFF5F7FA),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: Colors.grey.shade200),
              ),
              child: Text(
                queueNumber,
                style: const TextStyle(
                  fontSize: 32,
                  fontWeight: FontWeight.w800,
                  letterSpacing: 4,
                  color: Color(0xFF1A73E8),
                ),
              ),
            ),
            const SizedBox(height: 20),
            SizedBox(
              width: double.infinity,
              child: FilledButton(
                onPressed: () => Navigator.of(context).pop(),
                style: FilledButton.styleFrom(
                  backgroundColor: const Color(0xFF1A73E8),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                ),
                child: const Text('Tambah Lagi'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showSnackBar(
    String message, {
    bool isError = false,
    IconData icon = Icons.info_outline,
  }) {
    ScaffoldMessenger.of(context).clearSnackBars();
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Row(
          children: [
            Icon(icon, color: Colors.white, size: 18),
            const SizedBox(width: 10),
            Expanded(child: Text(message)),
          ],
        ),
        backgroundColor: isError ? Colors.red.shade700 : Colors.green.shade700,
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
        margin: const EdgeInsets.all(16),
        duration: const Duration(seconds: 3),
      ),
    );
  }

  // ─── UI ───────────────────────────────────────────────────────────────────
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF5F7FA),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        centerTitle: false,
        title: const Text(
          'Tambah Kendaraan ke Antrian',
          style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600),
        ),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(1),
          child: Divider(height: 1, color: Colors.grey.shade200),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // ── Plate ──────────────────────────────────────────────────────
              const _SectionLabel(label: 'Nomor Plat Kendaraan'),
              const SizedBox(height: 8),
              _PlateInputRow(
                controller: _plateController,
                onScanTap: scanPlate,
              ),
              const SizedBox(height: 8),
              Row(
                children: [
                  Icon(
                    Icons.document_scanner_outlined,
                    size: 13,
                    color: Colors.grey.shade400,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    'Ketik nomor plat secara manual.',
                    style: TextStyle(fontSize: 12, color: Colors.grey.shade400),
                  ),
                ],
              ),

              const SizedBox(height: 24),

              // ── Owner Info ─────────────────────────────────────────────────
              const _SectionLabel(label: 'Informasi Pemilik'),
              const SizedBox(height: 8),
              _OutlinedField(
                controller: _ownerNameController,
                hint: 'cth. Budi Santoso',
                label: 'Nama Lengkap',
                prefixIcon: Icons.person_outline,
                validator: (v) => (v == null || v.trim().isEmpty)
                    ? 'Nama tidak boleh kosong.'
                    : null,
              ),
              const SizedBox(height: 12),
              _OutlinedField(
                controller: _ownerPhoneController,
                hint: 'cth. 08123456789',
                label: 'Nomor Telepon',
                prefixIcon: Icons.phone_outlined,
                keyboardType: TextInputType.phone,
                inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                validator: (v) {
                  if (v == null || v.trim().isEmpty)
                    return 'Telepon tidak boleh kosong.';
                  if (v.trim().length < 9)
                    return 'Masukkan nomor telepon yang valid.';
                  return null;
                },
              ),

              const SizedBox(height: 24),

              // ── Vehicle Photo ──────────────────────────────────────────────
              const _SectionLabel(label: 'Foto Kendaraan'),
              const SizedBox(height: 4),
              Text(
                'Ambil foto kendaraan secara keseluruhan untuk catatan.',
                style: TextStyle(fontSize: 12, color: Colors.grey.shade500),
              ),
              const SizedBox(height: 10),
              _ImagePickerCard(image: _vehicleImage, onTap: pickVehicleImage),

              const SizedBox(height: 32),

              // ── Submit ─────────────────────────────────────────────────────
              _SubmitButton(isLoading: _isLoading, onPressed: submit),

              if (_isLoading) ...[
                const SizedBox(height: 12),
                const _SubmitStepHint(),
              ],
            ],
          ),
        ),
      ),
    );
  }
}

// ─── Sub-Widgets ──────────────────────────────────────────────────────────────

class _SectionLabel extends StatelessWidget {
  const _SectionLabel({required this.label, this.trailing});
  final String label;
  final Widget? trailing;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w600,
            color: Color(0xFF374151),
          ),
        ),
        if (trailing != null) ...[const SizedBox(width: 8), trailing!],
      ],
    );
  }
}

/// Plate input tanpa state [isScanning] karena OCR dihapus.
class _PlateInputRow extends StatelessWidget {
  const _PlateInputRow({required this.controller, required this.onScanTap});

  final TextEditingController controller;
  final VoidCallback onScanTap;

  @override
  Widget build(BuildContext context) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Expanded(
          child: TextFormField(
            controller: controller,
            textCapitalization: TextCapitalization.characters,
            style: const TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w700,
              letterSpacing: 2,
            ),
            decoration: InputDecoration(
              hintText: 'cth. B 1234 XYZ',
              filled: true,
              fillColor: Colors.white,
              prefixIcon: const Icon(Icons.directions_car_outlined),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              contentPadding: const EdgeInsets.symmetric(
                horizontal: 16,
                vertical: 14,
              ),
            ),
            validator: (v) {
              if (v == null || v.trim().isEmpty)
                return 'Plat tidak boleh kosong.';
              if (v.trim().length < 3) return 'Masukkan nomor plat yang valid.';
              return null;
            },
          ),
        ),
        const SizedBox(width: 10),
        // Tombol Scan tetap ditampilkan di UI namun memanggil stub
        SizedBox(
          height: 52,
          child: Tooltip(
            message: 'Fitur OCR belum tersedia',
            child: OutlinedButton(
              onPressed: onScanTap,
              style: OutlinedButton.styleFrom(
                foregroundColor: const Color(0xFF1A73E8),
                side: const BorderSide(color: Color(0xFF1A73E8), width: 1.5),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                padding: const EdgeInsets.symmetric(horizontal: 14),
              ),
              child: const Row(
                children: [
                  Icon(Icons.document_scanner_outlined, size: 18),
                  SizedBox(width: 6),
                  Text(
                    'Scan',
                    style: TextStyle(fontWeight: FontWeight.w600, fontSize: 14),
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class _OutlinedField extends StatelessWidget {
  const _OutlinedField({
    required this.controller,
    required this.hint,
    required this.label,
    required this.prefixIcon,
    this.keyboardType,
    this.inputFormatters,
    this.validator,
  });

  final TextEditingController controller;
  final String hint;
  final String label;
  final IconData prefixIcon;
  final TextInputType? keyboardType;
  final List<TextInputFormatter>? inputFormatters;
  final String? Function(String?)? validator;

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      controller: controller,
      keyboardType: keyboardType,
      inputFormatters: inputFormatters,
      decoration: InputDecoration(
        hintText: hint,
        labelText: label,
        filled: true,
        fillColor: Colors.white,
        prefixIcon: Icon(prefixIcon),
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 16,
          vertical: 14,
        ),
      ),
      validator: validator,
    );
  }
}

class _ImagePickerCard extends StatelessWidget {
  const _ImagePickerCard({required this.image, required this.onTap});
  final File? image;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        height: 220,
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(
            color: image != null
                ? const Color(0xFF1A73E8)
                : Colors.grey.shade300,
            width: image != null ? 2 : 1.5,
          ),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.04),
              blurRadius: 8,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        clipBehavior: Clip.hardEdge,
        child: image != null ? _buildPreview() : _buildPlaceholder(),
      ),
    );
  }

  Widget _buildPreview() {
    return Stack(
      fit: StackFit.expand,
      children: [
        Image.file(image!, fit: BoxFit.cover),
        Positioned(
          bottom: 0,
          left: 0,
          right: 0,
          child: Container(
            padding: const EdgeInsets.symmetric(vertical: 10),
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.bottomCenter,
                end: Alignment.topCenter,
                colors: [Colors.black.withOpacity(0.65), Colors.transparent],
              ),
            ),
            child: const Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.camera_alt, color: Colors.white, size: 16),
                SizedBox(width: 6),
                Text(
                  'Tap untuk ulangi',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 13,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildPlaceholder() {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Container(
          padding: const EdgeInsets.all(16),
          decoration: const BoxDecoration(
            color: Color(0xFFE8F0FE),
            shape: BoxShape.circle,
          ),
          child: const Icon(
            Icons.camera_alt_outlined,
            size: 32,
            color: Color(0xFF1A73E8),
          ),
        ),
        const SizedBox(height: 12),
        const Text(
          'Tap untuk ambil foto kendaraan',
          style: TextStyle(
            fontSize: 15,
            fontWeight: FontWeight.w500,
            color: Color(0xFF374151),
          ),
        ),
        const SizedBox(height: 4),
        Text(
          'Ambil foto kendaraan secara keseluruhan',
          style: TextStyle(fontSize: 13, color: Colors.grey.shade500),
        ),
      ],
    );
  }
}

class _SubmitButton extends StatelessWidget {
  const _SubmitButton({required this.isLoading, required this.onPressed});
  final bool isLoading;
  final VoidCallback onPressed;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 52,
      child: FilledButton(
        onPressed: isLoading ? null : onPressed,
        style: FilledButton.styleFrom(
          backgroundColor: const Color(0xFF1A73E8),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
        child: isLoading
            ? const SizedBox(
                width: 22,
                height: 22,
                child: CircularProgressIndicator(
                  color: Colors.white,
                  strokeWidth: 2.5,
                ),
              )
            : const Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.add_circle_outline, size: 20),
                  SizedBox(width: 8),
                  Text(
                    'Tambah ke Antrian',
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
                  ),
                ],
              ),
      ),
    );
  }
}

class _SubmitStepHint extends StatelessWidget {
  const _SubmitStepHint();

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Icon(
          Icons.cloud_upload_outlined,
          size: 14,
          color: Colors.grey.shade400,
        ),
        const SizedBox(width: 6),
        Text(
          'Mengunggah foto lalu membuat entri antrian...',
          style: TextStyle(fontSize: 12, color: Colors.grey.shade400),
        ),
      ],
    );
  }
}
