package com.sivalabs.bookmarks.api;


import com.sivalabs.bookmarks.AbstractIntegrationTest;
import com.sivalabs.bookmarks.domain.Bookmark;
import com.sivalabs.bookmarks.domain.BookmarkRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.CsvSource;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.context.jdbc.Sql;

import static io.restassured.RestAssured.given;
import static org.assertj.core.api.Assertions.assertThat;
import static org.hamcrest.Matchers.is;

@Sql("classpath:/test_data.sql")
class BookmarksControllerTests extends AbstractIntegrationTest {
    @Autowired
    BookmarkRepository repository;

    @Test
    void shouldFetchBookmarks() {
        given()
                .when()
                .get("/api/bookmarks")
                .then()
                .statusCode(200)
                .body("size()", is(6));
    }

    @Test
    void shouldGetNotFoundWhenBookmarkNotFound() {
        given()
                .when()
                .get("/api/bookmarks/9999")
                .then()
                .statusCode(404);
    }

    @Test
    void shouldCreateBookmark() {
        given()
                .contentType("application/json")
                .body("""
                    {
                        "title":"My Title",
                        "url":"https://sivalabs.in"
                    }
                """)
                .when()
                .post("/api/bookmarks")
                .then()
                .statusCode(201);
    }

    @ParameterizedTest
    @CsvSource(value= {
            "My Title,",
            ",https://sivalabs.in",
            ","
    })
    void shouldCreateBookmarkWithInvalidPayloadReturnBadRequest(String title, String url) {
        given()
                .contentType("application/json")
                .body("""
                    {
                        "title": %s,
                        "url": %s
                    }
                """.formatted(title, url))
                .when()
                .post("/api/bookmarks")
                .then()
                .statusCode(400);
    }

    @Test
    void updateBookmark() {
        given()
                .contentType("application/json")
                .body("""
                    {
                        "title":"My Blog",
                        "url":"https://sivalabs.in"
                    }
                """)
                .when()
                .put("/api/bookmarks/1")
                .then()
                .statusCode(200);

        Bookmark bookmark = repository.findById(1L).orElseThrow();
        assertThat(bookmark.title()).isEqualTo("My Blog");
        assertThat(bookmark.url()).isEqualTo("https://sivalabs.in");
    }

    @Test
    void shouldDeleteBookmark() {
        given()
                .when()
                .delete("/api/bookmarks/1")
                .then()
                .statusCode(204);

        assertThat(repository.findById(1L)).isEmpty();
    }
}
