import 'dart:convert';

import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:hive/hive.dart';
import 'package:http/http.dart';
import 'package:http/http.dart' as http;

Future<Response> get(String uri, [Map<String, String>? headers]) async {
  final serverUrl = dotenv.env['SERVER_URL'];
  if(serverUrl == null){
    throw Exception('serverUrl not found');
  }
  final url = Uri.parse(serverUrl + uri);
  print(url);
  return await http.get(url, headers: headers);
}

Future<Response> post(String uri, {Object? body, Map<String, String>? headers}) async {
  final serverUrl = dotenv.env['SERVER_URL'];
  if(serverUrl == null){
    throw Exception('serverUrl not found');
  }
  final url = Uri.parse(serverUrl + uri);
  print(url);
  return await http.post(
    url,
    body: jsonEncode(body),
    headers: headers
  );
}

Future<Response> authPost(String uri, {Object? body}) async {
  final userStorage = Hive.box('userStorage');
  final authToken = setAuthToken(userStorage.get("AuthToken").toString());
  return await post('auth/logout', body: body, headers: authToken);
}

Future<Response> authGet(String uri) async {
  final userStorage = Hive.box('userStorage');
  final authToken = setAuthToken(userStorage.get("AuthToken").toString());
  return await get(uri, authToken);
}


Map<String, String> setAuthToken(String token){
  return {
    "Authorization": "Bearer " + token
  };
}