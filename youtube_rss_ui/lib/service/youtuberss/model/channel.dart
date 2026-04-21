class Channel {
  final int id;
  final String url;
  final String name;
  final String rssUrl;

  Channel({
    required this.id,
    required this.url,
    required this.name,
    required this.rssUrl,
  });

  factory Channel.fromJson(Map<String, dynamic> json) {
    return Channel(
      id: json['id'] as int,
      url: json['url'] as String,
      name: json['name'] as String,
      rssUrl: json['rss_url'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'url': url,
      'name': name,
      'rss_url': rssUrl,
    };
  }
}
