import 'dart:io';

import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';

class CreateQueuePage extends StatefulWidget {
  const CreateQueuePage({super.key});

  @override
  State<CreateQueuePage> createState() => _CreateQueuePageState();
}

class _CreateQueuePageState extends State<CreateQueuePage> {
  // ─── State ────────────────────────────────────────────────────────────────
  final _plateController = TextEditingController();
  final _formKey = GlobalKey<FormState>();

  File? _vehicleImage;
  bool _isLoading = false;

  // ─── Lifecycle ────────────────────────────────────────────────────────────
  @override
  void dispose() {
    _plateController.dispose();
    super.dispose();
  }

  // ─── Logic ────────────────────────────────────────────────────────────────

  /// Opens the device camera, captures a photo, and stores it as a [File].
  /// Image quality is compressed to 60% to reduce upload size.
  Future<void> pickImage() async {
    final picker = ImagePicker();

    final XFile? picked = await picker.pickImage(
      source: ImageSource.camera,
      imageQuality: 60,
    );

    if (picked != null) {
      setState(() {
        _vehicleImage = File(picked.path);
      });
    }
  }

  /// Validates inputs, calls the service, and handles success/error feedback.
  Future<void> submit() async {
    // Validate plate field
    if (!_formKey.currentState!.validate()) return;

    // Validate image separately (not a form field)
    if (_vehicleImage == null) {
      _showSnackBar('Please take a photo of the vehicle.', isError: true);
      return;
    }

    setState(() => _isLoading = true);
  }

  void _resetForm() {
    _plateController.clear();
    setState(() => _vehicleImage = null);
  }

  void _showSnackBar(String message, {bool isError = false}) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: isError ? Colors.red.shade700 : Colors.green.shade700,
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
        margin: const EdgeInsets.all(16),
      ),
    );
  }

  // ─── UI ───────────────────────────────────────────────────────────────────
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: const Color(0xFFF5F7FA),
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        centerTitle: false,
        title: const Text(
          'Add Vehicle to Queue',
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
              _SectionLabel(label: 'Vehicle Plate Number'),
              const SizedBox(height: 8),
              _PlateTextField(controller: _plateController),
              const SizedBox(height: 28),
              _SectionLabel(label: 'Vehicle Photo'),
              const SizedBox(height: 8),
              _ImagePickerCard(image: _vehicleImage, onTap: pickImage),
              const SizedBox(height: 32),
              _SubmitButton(isLoading: _isLoading, onPressed: submit),
            ],
          ),
        ),
      ),
    );
  }
}

// ─── Sub-Widgets ──────────────────────────────────────────────────────────────
// Kept in the same file for MVP simplicity.
// Move to lib/features/queue/widgets/ when they grow or are reused elsewhere.

class _SectionLabel extends StatelessWidget {
  const _SectionLabel({required this.label});
  final String label;

  @override
  Widget build(BuildContext context) {
    return Text(
      label,
      style: const TextStyle(
        fontSize: 14,
        fontWeight: FontWeight.w600,
        color: Color(0xFF374151),
      ),
    );
  }
}

class _PlateTextField extends StatelessWidget {
  const _PlateTextField({required this.controller});
  final TextEditingController controller;

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      controller: controller,
      textCapitalization: TextCapitalization.characters,
      style: const TextStyle(
        fontSize: 18,
        fontWeight: FontWeight.w700,
        letterSpacing: 2,
      ),
      decoration: const InputDecoration(
        hintText: 'e.g. B 1234 XYZ',
        filled: true,
        fillColor: Colors.white,
        prefixIcon: Icon(Icons.directions_car_outlined),
      ),
      validator: (value) {
        if (value == null || value.trim().isEmpty) {
          return 'Plate number cannot be empty.';
        }
        if (value.trim().length < 3) {
          return 'Enter a valid plate number.';
        }
        return null;
      },
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
        // Retake overlay at bottom
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
                  'Tap to retake',
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
          decoration: BoxDecoration(
            color: const Color(0xFFE8F0FE),
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
          'Tap to open camera',
          style: TextStyle(
            fontSize: 15,
            fontWeight: FontWeight.w500,
            color: Color(0xFF374151),
          ),
        ),
        const SizedBox(height: 4),
        Text(
          'Take a clear photo of the vehicle',
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
                    'Add to Queue',
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
                  ),
                ],
              ),
      ),
    );
  }
}
