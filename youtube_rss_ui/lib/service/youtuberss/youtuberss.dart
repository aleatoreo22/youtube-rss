import 'dart:convert';
import 'package:http/http.dart' as http;
import 'model/content.dart';
import 'model/pagination.dart';
import 'model/user.dart';
import 'model/channel.dart';

class YoutubeRSS {
  static const String baseUrl = 'http://onehome:1234';

  Future<List<Content>> getContentByDate(DateTime date) async {
    final dateString =
        '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}';
    final response = await http.get(Uri.parse("$baseUrl/content/$dateString"));

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Content.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load content: ${response.statusCode}');
    }
  }

  Future<List<Content>> getUserContentByDate(int userId, DateTime date) async {
    final dateString =
        '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}';
    final response = await http.get(
      Uri.parse("$baseUrl/user/$userId/content/$dateString"),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => Content.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load user content: ${response.statusCode}');
    }
  }

  Future<void> addUserChannel(int userId, String channelUrl) async {
    final response = await http.post(
      Uri.parse("$baseUrl/user/$userId/channel"),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'url': channelUrl}),
    );

    if (response.statusCode != 200 && response.statusCode != 204) {
      throw Exception('Failed to add channel: ${response.statusCode}');
    }
  }

  Future<void> ping() async {
    final response = await http.get(Uri.parse("$baseUrl/ping"));

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      print(data);
    } else {
      throw Exception('Erro na requisição');
    }
  }

  Future<User> GetUser(int userId) async {
    final response = await http.get(
      Uri.parse("$baseUrl/user/$userId"),
    );

    if (response.statusCode == 200) {
      return User.fromJson(jsonDecode(response.body));
    } else if (response.statusCode == 404) {
      throw Exception('User not found');
    } else {
      throw Exception('Failed to get user: ${response.statusCode}');
    }
  }

  Future<List<Channel>> GetUserChannels(int userId) async {
    final response = await http.get(
      Uri.parse("$baseUrl/user/$userId/channel"),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body)['channels'];
      return data.map((json) => Channel.fromJson(json)).toList();
    } else if (response.statusCode == 404) {
      throw Exception('User not found');
    } else {
      throw Exception('Failed to get user channels: ${response.statusCode}');
    }
  }

  Future<void> DeleteUserChannel(int userId, int channelId) async {
    final response = await http.delete(
      Uri.parse("$baseUrl/user/$userId/channel/$channelId"),
    );

    if (response.statusCode == 200 || response.statusCode == 204) {
      // Success
    } else if (response.statusCode == 404) {
      throw Exception('User not found');
    } else {
      throw Exception('Failed to delete channel: ${response.statusCode}');
    }
  }

  Future<Pagination> getUserContentPaginated(
    int userId,
    int page,
    int limit,
  ) async {
    final response = await http.get(
      Uri.parse("$baseUrl/user/$userId/content?page=$page&limit=$limit"),
    );

    if (response.statusCode == 200) {
      return Pagination.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to load content: ${response.statusCode}');
    }
  }
}
