import { Component, OnInit, ViewChild } from '@angular/core';
import { Classify } from '../../model/config/classify';
import { ClassifyService } from '../../service/config/classify.service';
import { SyncClassifyComponent } from './sync/sync.component';

@Component({
  selector: 'app-classify',
  templateUrl: './classify.component.html',
  styleUrls: ['./classify.component.css']
})
export class ClassifyComponent implements OnInit {
  @ViewChild(SyncClassifyComponent)
  syncClassify: SyncClassifyComponent;

  classifys: Classify[];

  pageSize: number = 10;
  totalCount: number = 0;
  currentPage: number = 1;

  constructor(
    private classifyService: ClassifyService
  ) { }

  ngOnInit() {
    this.pageSize = 10;
    this.refresh();
  }

  syncd(syncd: boolean): void {
    if (syncd) {
      this.refresh();
    }
  }

  openSyncModel(id?: number): void {
    this.syncClassify.New(id);
  }

  delete(classify: Classify): void {
    this.classifyService.Delete(classify.Id)
    .subscribe(res => {
      this.refresh();
    })
  }

  load(state: any): void {
    if (state && state.page) {
      this.refreshClassify(state.page.from, state.page.to + 1);
    }
  }

  refresh() {
    this.currentPage = 1;
    this.refreshClassify(0, 10);
  }

  refreshClassify(from: number, to: number): void {
    this.classifyService.List()
    .subscribe(res => {
      this.totalCount = res.length;
      this.classifys = res.slice(from, to);
    })
  }
}
