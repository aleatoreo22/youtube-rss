class UserChannel {
  final int id;
  final int user;
  final int channel;

  UserChannel({
    required this.id,
    required this.user,
    required this.channel,
  });

  factory UserChannel.fromJson(Map<String, dynamic> json) {
    return UserChannel(
      id: json['id'] as int,
      user: json['user'] as int,
      channel: json['channel'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user': user,
      'channel': channel,
    };
  }
}
