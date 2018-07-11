import { Component, OnInit, Output, Input, EventEmitter } from '@angular/core'; 

import { Country }           from '../../../../model/time/country';
import { CountryService  }   from '../../../../service/time/country.service';

@Component({
  selector: 'time-country-save',
  templateUrl: './save.component.html',
  styleUrls: ['./save.component.css']
})

export class CountrySaveComponent implements OnInit {
  country: Country = new Country;
  modelOpened: boolean = false;

  @Output() save = new EventEmitter<boolean>();

  constructor(
    private countryService: CountryService,
  ) { }

  ngOnInit() {
  }

  New(id?: number): void {
    if (id) {
      this.countryService.Get(id).subscribe(res => {
        this.country = res;
        this.modelOpened = true;
      })
    } else {
      this.country = new Country();
      this.modelOpened = true;
    }
  }            

  Submit(): void {
    if (this.country.Id == null) {
      this.countryService.Add(this.country).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    } else {
      this.countryService.Update(this.country).subscribe(res => {
        this.modelOpened = false;
        this.save.emit(true);
      })
    }
  }
}
