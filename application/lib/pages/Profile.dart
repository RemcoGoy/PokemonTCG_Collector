import 'package:application/components/Button.dart';
import 'package:application/services/AuthService.dart';
import 'package:flutter/material.dart';

class Profile extends StatelessWidget {
  const Profile({super.key});

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Column(
        children: [
          Icon(
            Icons.person,
            size: 150,
          ),
          TCGButton(
            buttonText: 'Logout',
            onTap: () => logout(context),
          ),
        ]
      ),
    );
  }
}