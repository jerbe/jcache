package jcache

import (
	"context"

	"github.com/jerbe/jcache/v2/driver"
)

/**
  @author : Jerbe - The porter from Earth
  @time : 2023/9/28 12:43
  @describe :
*/

// SortedSetClient 已排序的集合客户端
type SortedSetClient struct {
	BaseClient
}

// NewSortedSetClient 返回一个已排序的集合客户端
func NewSortedSetClient(drivers ...driver.SortedSet) *SortedSetClient {
	drs := make([]driver.Common, 0)
	for i := 0; i < len(drivers); i++ {
		drs = append(drs, drivers[i])
	}

	if len(drs) == 0 {
		drs = append(drs, driver.NewMemory())
	}

	return &SortedSetClient{
		BaseClient: BaseClient{drivers: drs},
	}
}

// =======================================================
// ===================== SORTED SET ======================
// =======================================================

// ZAdd 添加有序集合的元素
func (cli *SortedSetClient) ZAdd(ctx context.Context, key string, members ...driver.Z) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for i, c := range cli.drivers {
		if v := c.(driver.SortedSet).ZAdd(ctx, key, members...); i == 0 {
			value = v
		}
	}
	return value
}

// ZCard 获取有序集合的元素数量
func (cli *SortedSetClient) ZCard(ctx context.Context, key string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for _, c := range cli.drivers {
		if value = c.(driver.SortedSet).ZCard(ctx, key); returnable(value) {
			return value
		}
	}
	return value
}

// ZCount 返回有序集 key 中， score 值在 min 和 max 之间(默认包括 score 值等于 min 或 max )的成员的数量。
func (cli *SortedSetClient) ZCount(ctx context.Context, key, min, max string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for _, c := range cli.drivers {
		if value = c.(driver.SortedSet).ZCount(ctx, key, min, max); returnable(value) {
			return value
		}
	}
	return value
}

// ZIncrBy 为有序集 key 的成员 member 的 score 值加上增量 increment 。
// 可以通过传递一个负数值 increment ，让 score 减去相应的值，比如 ZINCRBY key -5 member ，就是让 member 的 score 值减去 5
// @return member 成员的新 score 值
func (cli *SortedSetClient) ZIncrBy(ctx context.Context, key string, increment float64, member string) driver.FloatValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.FloatValuer
	for i, c := range cli.drivers {
		if v := c.(driver.SortedSet).ZIncrBy(ctx, key, increment, member); i == 0 {
			value = v
		}
	}
	return value
}

// ZRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递增(从小到大)来排序。
// 具有相同 score 值的成员按字典序(lexicographical order )来排列。
// 如果你需要成员按 score 值递减(从大到小)来排列，请使用 ZREVRANGE 命令。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (cli *SortedSetClient) ZRange(ctx context.Context, key string, start, stop int64) driver.StringSliceValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StringSliceValuer
	for _, c := range cli.drivers {
		if value = c.(driver.SortedSet).ZRange(ctx, key, start, stop); returnable(value) {
			return value
		}
	}
	return value
}

// ZRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。有序集成员按 score 值递增(从小到大)次序排列。
// 具有相同 score 值的成员按字典序(lexicographical order)来排列(该属性是有序集提供的，不需要额外的计算)。
// 可选的 LIMIT 参数指定返回结果的数量及区间(就像SQL中的 SELECT LIMIT offset, count )，注意当 offset 很大时，定位 offset 的操作可能需要遍历整个有序集，此过程最坏复杂度为 O(N) 时间。
func (cli *SortedSetClient) ZRangeByScore(ctx context.Context, key string, opt *driver.ZRangeBy) driver.StringSliceValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StringSliceValuer
	for _, c := range cli.drivers {
		if value = c.(driver.SortedSet).ZRangeByScore(ctx, key, opt); returnable(value) {
			return value
		}
	}
	return value
}

// ZRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递增(从小到大)顺序排列。
// 排名以 0 为底，也就是说， score 值最小的成员排名为 0 。
func (cli *SortedSetClient) ZRank(ctx context.Context, key, member string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for _, c := range cli.drivers {
		if value = c.(driver.SortedSet).ZRank(ctx, key, member); returnable(value) {
			return value
		}
	}
	return value
}

// ZRem 移除有序集 key 中的一个或多个成员，不存在的成员将被忽略。
// @return 被成功移除的成员的数量，不包括被忽略的成员
func (cli *SortedSetClient) ZRem(ctx context.Context, key string, members ...interface{}) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for i, c := range cli.drivers {
		if v := c.(driver.SortedSet).ZRem(ctx, key, members...); i == 0 {
			value = v
		}
	}
	return value
}

// ZRemRangeByRank 移除有序集 key 中，指定排名(rank)区间内的所有成员。
// 区间分别以下标参数 start 和 stop 指出，包含 start 和 stop 在内。
// 下标参数 start 和 stop 都以 0 为底，也就是说，以 0 表示有序集第一个成员，以 1 表示有序集第二个成员，以此类推。
// 你也可以使用负数下标，以 -1 表示最后一个成员， -2 表示倒数第二个成员，以此类推。
func (cli *SortedSetClient) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for i, c := range cli.drivers {
		if v := c.(driver.SortedSet).ZRemRangeByRank(ctx, key, start, stop); i == 0 {
			value = v
		}
	}
	return value
}

// ZRemRangeByScore 返回有序集 key 中，所有 score 值介于 min 和 max 之间(包括等于 min 或 max )的成员。
// 有序集成员按 score 值递增(从小到大)次序排列。
func (cli *SortedSetClient) ZRemRangeByScore(ctx context.Context, key, min, max string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for i, c := range cli.drivers {
		if v := c.(driver.SortedSet).ZRemRangeByScore(ctx, key, min, max); i == 0 {
			value = v
		}
	}
	return value
}

// ZRevRange 返回有序集 key 中，指定区间内的成员。
// 其中成员的位置按 score 值递减(从大到小)来排列。
// 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列。
func (cli *SortedSetClient) ZRevRange(ctx context.Context, key string, start, stop int64) driver.StringSliceValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.StringSliceValuer
	for _, c := range cli.drivers {
		if value = c.(driver.SortedSet).ZRevRange(ctx, key, start, stop); returnable(value) {
			return value
		}
	}
	return value
}

// ZRevRank 返回有序集 key 中成员 member 的排名。其中有序集成员按 score 值递减(从大到小)排序。
// 排名以 0 为底，也就是说， score 值最大的成员排名为 0 。
func (cli *SortedSetClient) ZRevRank(ctx context.Context, key, member string) driver.IntValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.IntValuer
	for _, c := range cli.drivers {
		if value = c.(driver.SortedSet).ZRevRank(ctx, key, member); returnable(value) {
			return value
		}
	}
	return value
}

// ZScore 返回有序集 key 中，成员 member 的 score 值。
// 如果 member 元素不是有序集 key 的成员，或 key 不存在，返回 nil 。
func (cli *SortedSetClient) ZScore(ctx context.Context, key, member string) driver.FloatValuer {
	ctx, _ = cli.preCheck(ctx)

	var value driver.FloatValuer
	for _, c := range cli.drivers {
		if value = c.(driver.SortedSet).ZScore(ctx, key, member); returnable(value) {
			return value
		}
	}
	return value
}
