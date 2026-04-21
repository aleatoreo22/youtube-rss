import 'content.dart';

class Pagination {
  final int page;
  final int limit;
  final int total;
  final int totalPages;
  final List<Content> items;

  Pagination({
    required this.page,
    required this.limit,
    required this.total,
    required this.totalPages,
    required this.items,
  });

  factory Pagination.fromJson(Map<String, dynamic> json) {
    return Pagination(
      page: json['page'] ?? 1,
      limit: json['limit'] ?? 20,
      total: json['total'] ?? 0,
      totalPages: json['totalPages'] ?? 0,
      items: json['items'] != null
          ? (json['items'] as List).map((e) => Content.fromJson(e)).toList()
          : [],
    );
  }
}
