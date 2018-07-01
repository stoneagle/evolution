import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute }               from "@angular/router";

import { Item }                 from '../../../../model/quant/item';
import { Pool }                 from '../../../../model/quant/pool';
import { Classify }             from '../../../../model/quant/classify';
import { PoolService }          from '../../../../service/quant/pool.service';
import { ItemService }          from '../../../../service/quant/item.service';
import { PoolAddItemComponent } from '../add-item/pool-add-item.component';

@Component({
  selector: 'pool-list-item',
  templateUrl: './pool-list-item.component.html',
  styleUrls: ['./pool-list-item.component.css']
})
export class PoolListItemComponent implements OnInit {
  @ViewChild(PoolAddItemComponent)
  poolAddItem: PoolAddItemComponent;

  pool: Pool;
  items: Item[] = []; 
  selectedDeleteItems: Item[] = [];
  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private route: ActivatedRoute,
    private poolService: PoolService,
    private itemService: ItemService
  ) { 
    this.route.params.subscribe( params => {
      this.pool = new Pool();
      this.pool.Item = [];
      this.pool.Id = +params.id;
    });
  }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  added(added: boolean): void {
    if (added) {
      this.refresh();
    }
  }

  openAddModel(): void {
    this.poolAddItem.New(this.pool);
  }

  onDelete(): void {
    if (this.selectedDeleteItems.length > 0) {
      let deletePool = this.pool;
      deletePool.Item = this.selectedDeleteItems;
      this.poolService.DeleteItems(deletePool).subscribe(res => {
        if (res) {
          this.refresh();
        }
      })
    }
  }

  load(state: any): void {
    if (state && state.page) {
      this.items = this.pool.Item.slice(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshItem(0, 10);
  }

  refreshItem(from: number, to: number): void {
    this.poolService.Get(this.pool.Id).subscribe(res => {
      this.pool = res;
      this.totalCount = res.Item.length;
      this.items = res.Item.slice(from, to);
    })
  }

  syncPoint(item: Item): void {
    this.itemService.SyncPoint(item.Id).subscribe(res => {
    })
  }
}
