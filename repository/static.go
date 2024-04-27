package repository

import (
	"context"
	"github.com/MuxiKeStack/be-static/domain"
	"github.com/MuxiKeStack/be-static/pkg/logger"
	"github.com/MuxiKeStack/be-static/repository/cache"
	"github.com/MuxiKeStack/be-static/repository/dao"
	"time"
)

type StaticRepository interface {
	GetStaticByName(ctx context.Context, name string) (domain.Static, error)
	SaveStatic(ctx context.Context, static domain.Static) error
}

type CachedStaticRepository struct {
	dao   dao.StaticDAO
	cache cache.StaticCache
	l     logger.Logger
}

func NewCachedStaticRepository(dao dao.StaticDAO, cache cache.StaticCache, l logger.Logger) StaticRepository {
	return &CachedStaticRepository{dao: dao, cache: cache, l: l}
}

func (repo *CachedStaticRepository) GetStaticByName(ctx context.Context, name string) (domain.Static, error) {
	res, err := repo.cache.GetStatic(ctx, name)
	if err == nil {
		return res, nil
	}
	static, err := repo.dao.GetStaticByName(ctx, name)
	res = domain.Static{
		Name:    static.Name,
		Content: static.Content,
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := repo.cache.SetStatic(ctx, res)
		repo.l.Error("回写静态资源失败", logger.Error(er))
	}()
	return res, err
}

func (repo *CachedStaticRepository) SaveStatic(ctx context.Context, static domain.Static) error {
	err := repo.dao.Upsert(ctx, dao.Static{
		Name:    static.Name,
		Content: static.Content,
	})
	if err != nil {
		return err
	}
	// 更新缓存，或者立刻设置缓存
	return repo.cache.SetStatic(ctx, static)
}
