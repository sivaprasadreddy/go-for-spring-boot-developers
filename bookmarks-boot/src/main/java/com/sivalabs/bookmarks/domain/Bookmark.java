package com.sivalabs.bookmarks.domain;

import jakarta.validation.constraints.NotBlank;

import java.time.Instant;

public record Bookmark(
        Long id,
        @NotBlank(message = "Title is required")
        String title,
        @NotBlank(message = "URL is required")
        String url,
        Instant createdAt) {
}