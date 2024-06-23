create view if not exists Indexes(id, keyword)
as
select
    CK.comic_id,
    K.word
from
    ComicsKeywords CK
join Keywords K on K.id = CK.keyword_id;