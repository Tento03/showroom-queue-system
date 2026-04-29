import 'package:flutter/material.dart';
import 'package:vehicle_queue_app/features/queue/pages/create_queue_page.dart';
import 'package:vehicle_queue_app/features/queue/pages/queue_list_page.dart';

class MainNavigationPage extends StatefulWidget {
  const MainNavigationPage({super.key});

  @override
  State<MainNavigationPage> createState() => _MainNavigationPageState();
}

class _MainNavigationPageState extends State<MainNavigationPage> {
  int _currentIndex = 0;

  // Gunakan IndexedStack agar state tiap halaman dipertahankan
  final List<Widget> _pages = const [CreateQueuePage(), QueueListPage()];

  void _onItemTapped(int index) {
    setState(() => _currentIndex = index);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(index: _currentIndex, children: _pages),
      bottomNavigationBar: _BottomNav(
        currentIndex: _currentIndex,
        onTap: _onItemTapped,
      ),
    );
  }
}

// ─── Custom Bottom Navigation ─────────────────────────────────────────────────

class _BottomNav extends StatelessWidget {
  const _BottomNav({required this.currentIndex, required this.onTap});
  final int currentIndex;
  final ValueChanged<int> onTap;

  static const _items = [
    _NavItem(
      icon: Icons.add_circle_outline_rounded,
      activeIcon: Icons.add_circle_rounded,
      label: 'Generate',
    ),
    _NavItem(
      icon: Icons.format_list_bulleted_rounded,
      activeIcon: Icons.format_list_bulleted_rounded,
      label: 'Daftar Antrian',
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.08),
            blurRadius: 20,
            offset: const Offset(0, -4),
          ),
        ],
      ),
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
          child: Row(
            children: List.generate(_items.length, (i) {
              return Expanded(
                child: _NavButton(
                  item: _items[i],
                  isActive: currentIndex == i,
                  onTap: () => onTap(i),
                ),
              );
            }),
          ),
        ),
      ),
    );
  }
}

class _NavItem {
  const _NavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
  });
  final IconData icon;
  final IconData activeIcon;
  final String label;
}

class _NavButton extends StatelessWidget {
  const _NavButton({
    required this.item,
    required this.isActive,
    required this.onTap,
  });

  final _NavItem item;
  final bool isActive;
  final VoidCallback onTap;

  static const _activeColor = Color(0xFF1A73E8);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      behavior: HitTestBehavior.opaque,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        curve: Curves.easeInOut,
        padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 4),
        decoration: BoxDecoration(
          color: isActive ? _activeColor.withOpacity(0.08) : Colors.transparent,
          borderRadius: BorderRadius.circular(14),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            // ── Icon ──────────────────────────────────────────────────────────
            AnimatedSwitcher(
              duration: const Duration(milliseconds: 180),
              transitionBuilder: (child, anim) =>
                  ScaleTransition(scale: anim, child: child),
              child: Icon(
                isActive ? item.activeIcon : item.icon,
                key: ValueKey(isActive),
                color: isActive ? _activeColor : Colors.grey.shade400,
                size: 26,
              ),
            ),
            const SizedBox(height: 4),
            // ── Label ─────────────────────────────────────────────────────────
            AnimatedDefaultTextStyle(
              duration: const Duration(milliseconds: 200),
              style: TextStyle(
                fontSize: 11,
                fontWeight: isActive ? FontWeight.w700 : FontWeight.w500,
                color: isActive ? _activeColor : Colors.grey.shade400,
              ),
              child: Text(item.label, textAlign: TextAlign.center),
            ),

            // ── Active Dot Indicator ──────────────────────────────────────────
            const SizedBox(height: 4),
            AnimatedContainer(
              duration: const Duration(milliseconds: 200),
              width: isActive ? 6 : 0,
              height: isActive ? 6 : 0,
              decoration: const BoxDecoration(
                color: _activeColor,
                shape: BoxShape.circle,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
